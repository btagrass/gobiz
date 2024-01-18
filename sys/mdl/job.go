package mdl

import (
	"github.com/btagrass/gobiz/mdl"
)

type Job struct {
	mdl.Mdl
	Name     string `gorm:"uniqueIndex;not null" json:"name"`
	Desc     string `gorm:"" json:"desc"`
	Cron     string `gorm:"not null" json:"cron"`
	Arg      string `gorm:"" json:"arg"`
	ArgDesc  string `gorm:"" json:"argDesc"`
	Instance int    `gorm:"not null" json:"instance"`
	State    int    `gorm:"not null" json:"state"`
}

func (m Job) TableName() string {
	return "sys_job"
}
