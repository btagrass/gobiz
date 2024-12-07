package mdl

import (
	"github.com/btagrass/gobiz/mdl"
)

type Visit struct {
	mdl.Mdl
	UserId    int64  `gorm:"not null" json:"userId"`
	UserName  string `gorm:"not null" json:"userName"`
	Ip        string `gorm:"not null" json:"ip"`
	Method    string `gorm:"not null" json:"method"`
	Url       string `gorm:"not null" json:"url"`
	UserAgent string `gorm:"not null" json:"userAgent"`
	Req       string `gorm:"size:1000;not null" json:"req"`
	Resp      string `gorm:"size:1000;not null" json:"resp"`
}

func (m Visit) TableName() string {
	return "sys_visit"
}
