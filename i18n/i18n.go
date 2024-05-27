package i18n

import (
	"fmt"
	"io"

	"github.com/spf13/viper"
)

var (
	locale *viper.Viper
)

func Load(file io.Reader, types ...string) error {
	locale = viper.New()
	if len(types) == 0 {
		types = []string{"yaml"}
	}
	for _, t := range types {
		locale.SetConfigType(t)
	}
	return locale.ReadConfig(file)
}

func T(key string, args ...any) string {
	format := locale.GetString(key)
	return fmt.Sprintf(format, args...)
}
