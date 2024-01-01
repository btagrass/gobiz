package trc

import (
	"context"
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

const (
	StateCanceled = "Canceled"
	StateStarted  = "Started"
)

type Trc[T ITask] struct {
	*svc.Svc[T]
	hostname   string
	lastTasks  map[string]T
	timeshares [24]int
}

func NewTrc[T ITask](prefix string, timeshares [24]int) *Trc[T] {
	hostname, _ := os.Hostname()
	t := &Trc[T]{
		Svc:        svc.NewSvc[T](prefix),
		timeshares: timeshares,
		hostname:   hostname,
	}
	stateKey := fmt.Sprintf("%s:*:state", prefix)
	stateKeys := t.Redis.Keys(context.Background(), stateKey).Val()
	for _, k := range stateKeys {
		v := t.Redis.Get(context.Background(), k).Val()
		if strings.HasPrefix(v, hostname) {
			t.Redis.Del(context.Background(), k)
		}
	}
	return t
}

func (t *Trc[T]) Clean(duration time.Duration) error {
	currentTime := time.Now()
	keys, err := t.Redis.Keys(context.Background(), t.GetFullKey("*", "timeshares")).Result()
	if err != nil {
		return err
	}
	for _, k := range keys {
		ks := utl.Split(k, ':')
		if len(ks) < 3 {
			continue
		}
		dateTime, err := time.Parse("20060102", ks[2])
		if err != nil {
			logrus.Error(err)
			continue
		}
		if currentTime.Sub(dateTime) > duration {
			t.Redis.Del(context.Background(), k)
		}
	}
	return nil
}

func (t *Trc[T]) GetState(task T) string {
	stateKey := t.GetFullKey(task.GetCode(), "state")
	stateValue := t.Redis.Get(context.Background(), stateKey).Val()
	return strings.TrimPrefix(stateValue, t.hostname)
}

func (t *Trc[T]) Run(tasks []T, process func(T) error) error {
	// Stop
	thisTasks := make(map[string]T)
	for _, task := range tasks {
		thisTasks[task.GetCode()] = task
	}
	for _, lt := range t.lastTasks {
		_, ok := thisTasks[lt.GetCode()]
		if !ok {
			t.cancel(lt)
		}
	}
	// Match
	matchedTasks := make(map[string]T)
	for _, tt := range thisTasks {
		lt, ok := t.lastTasks[tt.GetCode()]
		if ok && (reflect.DeepEqual(lt, tt) || t.cancel(lt)) {
			continue
		}
		if t.isAvailable(tt) {
			matchedTasks[tt.GetCode()] = tt
		}
	}
	// Run
	for _, mt := range matchedTasks {
		go func(task T) {
			started := t.start(task)
			logrus.WithField("code", task.GetCode()).Debug("Start ", started, ". ", task)
			if started {
				err := process(task)
				if err != nil {
					logrus.WithField("code", task.GetCode()).Error(err)
				}
				stopped := t.stop(task)
				logrus.WithField("code", task.GetCode()).Debug("Stop ", stopped, ". ", task)
			}
		}(mt)
	}
	// Archive
	t.lastTasks = matchedTasks
	return nil
}

func (t *Trc[T]) calcPercent(beginTime, endTime time.Time) float64 {
	percent := 0.0
	timeshares := append(t.timeshares[:], t.timeshares[:]...)
	hourCount := cast.ToInt(math.Ceil(endTime.Sub(beginTime).Hours()))
	if hourCount < 0 {
		return percent
	}
	hours := timeshares[beginTime.Hour() : beginTime.Hour()+hourCount]
	beginMinute, endMinute := float64(beginTime.Minute()), float64(endTime.Minute())
	for i, j := 0, hourCount-1; i <= j; i++ {
		hourPercent := float64(hours[i])
		if i == 0 {
			if i < j {
				percent += hourPercent * (60 - beginMinute) / 60
			} else {
				percent += hourPercent * (endMinute - beginMinute) / 60
			}
		} else if i < j {
			percent += hourPercent
		} else {
			percent += hourPercent * endMinute / 60
		}
	}
	return percent
}

func (t *Trc[T]) calcRate(beginTime, currentTime, endTime time.Time) float64 {
	currentPercent := t.calcPercent(beginTime, currentTime)
	totalPercent := t.calcPercent(beginTime, endTime)
	return cast.ToFloat64(fmt.Sprintf("%.4f", currentPercent/totalPercent))
}

func (t *Trc[T]) cancel(task T) bool {
	stateKey := t.GetFullKey(task.GetCode(), "state")
	stateValue := fmt.Sprintf("%s.%s", t.hostname, StateCanceled)
	return t.Redis.SetXX(context.Background(), stateKey, stateValue, time.Hour).Val()
}

func (t *Trc[T]) isAvailable(task T) bool {
	expectedRate := t.calcRate(task.GetBeginTime(), time.Now(), task.GetEndTime())
	expectedCount := cast.ToInt64(cast.ToFloat64(task.GetCount()) * expectedRate)
	timesharesKey := t.GetFullKey(task.GetBeginTime().Format("20060102"), task.GetCode(), "timeshares")
	expectedKey := "expected"
	t.Redis.HSet(context.Background(), timesharesKey, expectedKey, expectedCount)
	actualKey := "actual"
	actualCount := cast.ToInt64(t.Redis.HGet(context.Background(), timesharesKey, actualKey).Val())
	return actualCount < expectedCount
}

func (t *Trc[T]) start(task T) bool {
	stateKey := t.GetFullKey(task.GetCode(), "state")
	stateValue := fmt.Sprintf("%s.%s", t.hostname, StateStarted)
	if !t.Redis.SetNX(context.Background(), stateKey, stateValue, time.Hour).Val() {
		return false
	}
	timesharesKey := t.GetFullKey(task.GetBeginTime().Format("20060102"), task.GetCode(), "timeshares")
	expectedKey := "expected"
	expectedCount := cast.ToInt64(t.Redis.HGet(context.Background(), timesharesKey, expectedKey).Val())
	actualKey := "actual"
	if t.Redis.HIncrBy(context.Background(), timesharesKey, actualKey, 1).Val() > expectedCount {
		t.Redis.HIncrBy(context.Background(), timesharesKey, actualKey, -1)
		return false
	}
	return true
}

func (t *Trc[T]) stop(task T) bool {
	timesharesKey := t.GetFullKey(task.GetBeginTime().Format("20060102"), task.GetCode(), "timeshares")
	state := t.GetState(task)
	if state == StateCanceled {
		state = "canceled"
	} else {
		state = "finished"
	}
	t.Redis.HIncrBy(context.Background(), timesharesKey, state, 1)
	stateKey := t.GetFullKey(task.GetCode(), "state")
	t.Redis.Del(context.Background(), stateKey).Val()
	return true
}
