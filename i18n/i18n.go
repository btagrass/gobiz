package i18n

import (
	"fmt"
	"io"

	"github.com/btagrass/gobiz/utl"
	"github.com/spf13/cast"
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
	for i := 0; i < len(args); i++ {
		keys := utl.Split(cast.ToString(args[i]), '$')
		var val string
		for _, k := range keys {
			val += locale.GetString(k)
		}
		args[i] = val
	}
	return fmt.Sprintf(format, args...)
}
