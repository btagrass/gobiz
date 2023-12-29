package mdl

import "github.com/btagrass/gobiz/mdl"

type Dict struct {
	mdl.Mdl
	Type     string `gorm:"not null" json:"type"`
	Code     int8   `gorm:"not null" json:"code"`
	Name     string `gorm:"not null" json:"name"`
	Sequence int    `gorm:"not null" json:"sequence"`
}

func (m Dict) TableName() string {
	return "sys_dict"
}
