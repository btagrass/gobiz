package utl

import (
	"log/slog"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/samber/lo"
)

func DifferenceFunc[T any](s1, s2 []T, key func(e T) string) ([]T, []T) {
	m1 := lo.Associate(s1, func(item T) (string, T) {
		return key(item), item
	})
	m2 := lo.Associate(s2, func(item T) (string, T) {
		return key(item), item
	})
	for k := range m1 {
		_, ok := m2[k]
		if ok {
			delete(m1, k)
			delete(m2, k)
		}
	}
	return lo.Values(m1), lo.Values(m2)
}

func ForParallel[T any](s []T, iterate func(e T) error, callback func(i int), size int) {
	var group sync.WaitGroup
	var locker sync.Mutex
	var index int
	pool, err := ants.NewPoolWithFunc(size, func(e any) {
		err := iterate(e.(T))
		if err != nil {
			slog.Error(err.Error())
		}
		locker.Lock()
		index++
		locker.Unlock()
		callback(index)
		group.Done()
	})
	if err != nil {
		slog.Error(err.Error())
	}
	defer pool.Release()
	for _, e := range s {
		group.Add(1)
		err = pool.Invoke(e)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	group.Wait()
}
