package mdl

import "github.com/btagrass/gobiz/mdl"

type Resource struct {
	mdl.Mdl
	ParentId int64      `gorm:"not null" json:"parentId"`
	Name     string     `gorm:"not null" json:"name"`
	Type     int8       `gorm:"not null" json:"type"`
	Icon     string     `gorm:"" json:"icon"`
	Uri      string     `gorm:"not null" json:"uri"`
	Act      string     `gorm:"" json:"act"`
	Sequence int        `gorm:"not null" json:"sequence"`
	Children []Resource `gorm:"foreignKey:ParentId" json:"children"`
}

func (m Resource) TableName() string {
	return "sys_resource"
}
