package mdl

import "github.com/btagrass/gobiz/mdl"

type User struct {
	mdl.Mdl
	DeptId   int64  `gorm:"not null" json:"deptId"`
	UserName string `gorm:"uniqueIndex;size:60;not null" json:"userName"`
	FullName string `gorm:"" json:"fullName"`
	Mobile   string `gorm:"not null" json:"mobile"`
	Password string `gorm:"size:60;not null" json:"password"`
	Frozen   bool   `gorm:"not null" json:"frozen"`
	Token    string `gorm:"-" json:"token"`
	Dept     *Dept  `gorm:"" json:"dept"`
}

func (m User) TableName() string {
	return "sys_user"
}
