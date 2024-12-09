package utl

import (
	"log/slog"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/samber/lo"
)

func Difference[T any](s1, s2 []T, equal func(t1, t2 T) bool) ([]T, []T) {
	var s1d, s2d []T
	for _, t1 := range s1 {
		if !lo.ContainsBy(s2, func(t2 T) bool {
			return equal(t1, t2)
		}) {
			s1d = append(s1d, t1)
		}
	}
	for _, t2 := range s2 {
		if !lo.ContainsBy(s1, func(t1 T) bool {
			return equal(t1, t2)
		}) {
			s2d = append(s2d, t2)
		}
	}
	return s1d, s2d
}

func ForParallel[T any](s []T, iterate func(t T) error, callback func(i int), size int) {
	var group sync.WaitGroup
	var locker sync.Mutex
	var finished int
	pool, err := ants.NewPoolWithFunc(size, func(a any) {
		err := iterate(a.(T))
		if err != nil {
			slog.Error(err.Error())
		}
		group.Done()
		locker.Lock()
		finished++
		locker.Unlock()
		callback(finished)
	})
	if err != nil {
		slog.Error(err.Error())
	}
	defer pool.Release()
	for _, v := range s {
		group.Add(1)
		err = pool.Invoke(v)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	group.Wait()
}

func Intersect[T any](s1, s2 []T, equal func(t1, t2 T) bool) []T {
	var s1i []T
	for _, t1 := range s1 {
		if lo.ContainsBy(s2, func(t2 T) bool {
			return equal(t1, t2)
		}) {
			s1i = append(s1i, t1)
		}
	}
	return s1i
}
