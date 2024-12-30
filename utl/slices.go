package utl

import (
	"log/slog"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/samber/lo"
)

func DifferenceFunc[T1, T2 any](s1 []T1, s2 []T2, k1 func(e1 T1) string, k2 func(e2 T2) string) ([]T1, []T2) {
	m1 := lo.Associate(s1, func(item T1) (string, T1) {
		return k1(item), item
	})
	m2 := lo.Associate(s2, func(item T2) (string, T2) {
		return k2(item), item
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
		if callback != nil {
			callback(index)
		}
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
