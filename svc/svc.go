package svc

import (
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/samber/do"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var (
	injector = do.New()
)

func Inject[T any](provider do.Provider[T]) {
	do.Provide(injector, provider)
	_ = do.MustInvoke[T](injector)
}

func Use[T any]() T {
	return do.MustInvoke[T](injector)
}

type Svc[M any] struct {
	Local  *cache.Cache
	Redis  *redis.Client
	prefix string
}

func NewSvc[M any](prefix string) *Svc[M] {
	s := &Svc[M]{
		Local: cache.New(cache.NoExpiration, 5*time.Minute),
	}
	addr := viper.GetString("redis.addr")
	if addr != "" {
		s.Redis = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: viper.GetString("redis.password"),
			DB:       viper.GetInt("redis.db"),
		})
	}
	s.prefix = prefix
	return s
}

func (s *Svc[M]) GetFullKey(keys ...any) string {
	var builder strings.Builder
	builder.WriteString(s.prefix)
	for _, k := range keys {
		builder.WriteString(":")
		builder.WriteString(cast.ToString(k))
	}
	return builder.String()
}
