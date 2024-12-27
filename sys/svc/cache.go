package svc

import (
	"context"
	"fmt"
	"time"

	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/utl"
	"github.com/bytedance/sonic"
	"github.com/samber/do"
	"github.com/spf13/cast"
)

type CacheSvc struct {
	*svc.Svc[mdl.Cache]
}

func NewCacheSvc(i *do.Injector) (*CacheSvc, error) {
	return &CacheSvc{
		Svc: svc.NewSvc[mdl.Cache]("sys:caches"),
	}, nil
}

func (s *CacheSvc) ListCaches(typ string, keywords ...string) ([]mdl.Cache, int64, error) {
	caches := make([]mdl.Cache, 0)
	if typ == "1" {
		i := 0
		for k, v := range s.Local.Items() {
			if utl.Contains(k, keywords...) {
				valString, _ := sonic.MarshalString(v.Object)
				var e int64
				if v.Expiration > 0 {
					e = cast.ToInt64(time.Until(time.Unix(0, v.Expiration)).Seconds())
				}
				caches = append(caches, mdl.Cache{
					Id:         cast.ToInt64(i + 1),
					Key:        k,
					Val:        v.Object,
					ValString:  valString,
					Expiration: e,
				})
				i++
			}
		}
	} else if typ == "2" {
		var keys []string
		for _, k := range keywords {
			keys = append(keys, s.Redis.Keys(context.Background(), fmt.Sprintf("*%s*", k)).Val()...)
		}
		for i, k := range keys {
			var v any
			t := s.Redis.Type(context.Background(), k).Val()
			if t == "hash" {
				v = s.Redis.HGetAll(context.Background(), k).Val()
			} else {
				v = s.Redis.Get(context.Background(), k).Val()
			}
			valString, _ := sonic.MarshalString(v)
			e := cast.ToInt64(s.Redis.TTL(context.Background(), k).Val().Seconds())
			caches = append(caches, mdl.Cache{
				Id:         cast.ToInt64(i + 1),
				Key:        k,
				Val:        v,
				ValString:  valString,
				Expiration: e,
			})
		}
	}
	return caches, cast.ToInt64(len(caches)), nil
}

func (s *CacheSvc) RemoveCaches(typ string, keys ...string) error {
	if typ == "local" {
		for _, k := range keys {
			s.Local.Delete(k)
		}
	} else if typ == "redis" {
		return s.Redis.Del(context.Background(), keys...).Err()
	}
	return nil
}
