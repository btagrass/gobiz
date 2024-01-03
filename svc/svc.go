package svc

import (
	"strings"

	"github.com/btagrass/gobiz/svc/internal"
	"github.com/samber/do"
	"github.com/spf13/cast"
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
	*internal.Cache
	prefix string
}

func NewSvc[M any](prefix string) *Svc[M] {
	return &Svc[M]{
		Cache:  internal.NewCache(),
		prefix: prefix,
	}
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
