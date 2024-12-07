package mdl

import (
	"github.com/btagrass/gobiz/mdl"
)

type Dept struct {
	mdl.Mdl
	ParentId int64  `gorm:"not null" json:"parentId"`
	Name     string `gorm:"not null" json:"name"`
	Phone    string `gorm:"" json:"phone"`
	Addr     string `gorm:"" json:"addr"`
	Sequence int    `gorm:"not null" json:"sequence"`
	Children []Dept `gorm:"foreignKey:ParentId" json:"children"`
}

func (m Dept) TableName() string {
	return "sys_dept"
}
