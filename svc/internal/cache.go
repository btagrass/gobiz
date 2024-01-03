package internal

import (
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
)

var (
	cch  *Cache
	once sync.Once
)

type Cache struct {
	Local *cache.Cache
	Redis *redis.Client
}

func NewCache() *Cache {
	once.Do(func() {
		cch = &Cache{
			Local: cache.New(cache.NoExpiration, 5*time.Minute),
		}
		addr := viper.GetString("redis.addr")
		if addr != "" {
			cch.Redis = redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: viper.GetString("redis.password"),
				DB:       viper.GetInt("redis.db"),
			})
		}
	})
	return cch
}
