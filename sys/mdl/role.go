package mdl

import "github.com/btagrass/gobiz/mdl"

type Role struct {
	mdl.Mdl
	Name string `gorm:"not null" json:"name"`
}

func (m Role) TableName() string {
	return "sys_role"
}
