package types

import (
	"database/sql/driver"
	"strings"

	"github.com/btagrass/gobiz/utl"
	"github.com/spf13/cast"
)

type StringSlice []string

func (s *StringSlice) Scan(value any) error {
	v, err := cast.ToStringE(value)
	if err != nil {
		return err
	}
	*s = utl.Split(v, ',')
	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	return strings.Join(s, ","), nil
}
