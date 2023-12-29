package svc

import (
	"context"
	"time"

	"github.com/btagrass/gobiz/dao"
)

type DataSvc[M any] struct {
	*Svc[M]
	*dao.Dao[M]
}

func NewDataSvc[M any](prefix string) *DataSvc[M] {
	return &DataSvc[M]{
		Svc: NewSvc[M](prefix),
		Dao: dao.NewDao[M](),
	}
}

func (s *DataSvc[M]) GetLocal(expiration time.Duration, conds ...any) (*M, error) {
	var m *M
	key := s.GetFullKey(conds)
	v, ok := s.Local.Get(key)
	if ok {
		m = v.(*M)
	} else {
		var err error
		m, err = s.Get(conds...)
		if err != nil {
			return m, err
		}
		s.Local.Set(key, m, expiration)
	}
	return m, nil
}

func (s *DataSvc[M]) ListLocal(expiration time.Duration, conds ...any) ([]M, error) {
	var ms []M
	key := s.GetFullKey(conds)
	v, ok := s.Local.Get(key)
	if ok {
		ms = v.([]M)
	} else {
		var err error
		ms, _, err = s.List(conds)
		if err != nil {
			return ms, err
		}
		s.Local.Set(key, ms, expiration)
	}
	return ms, nil
}

func (s *DataSvc[M]) GetRedis(expiration time.Duration, conds ...any) (*M, error) {
	var m *M
	key := s.GetFullKey(conds)
	err := s.Redis.Get(context.Background(), key).Scan(&m)
	if err != nil {
		m, err = s.Get(conds...)
		if err != nil {
			return m, err
		}
		err = s.Redis.Set(context.Background(), key, m, expiration).Err()
		if err != nil {
			return m, err
		}
	}
	return m, nil
}

func (s *DataSvc[M]) ListRedis(expiration time.Duration, conds ...any) ([]M, error) {
	var ms []M
	key := s.GetFullKey(conds)
	err := s.Redis.Get(context.Background(), key).Scan(&ms)
	if err != nil {
		ms, _, err = s.List(conds)
		if err != nil {
			return ms, err
		}
		err = s.Redis.Set(context.Background(), key, ms, expiration).Err()
		if err != nil {
			return ms, err
		}
	}
	return ms, nil
}
