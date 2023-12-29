package mdl

import (
	"fmt"
	"time"

	"github.com/bytedance/sonic"
	"gorm.io/gorm"
)

type Mdl struct {
	Id        int64          `gorm:"primaryKey;autoIncrement:false" json:"id"`
	CreatedAt time.Time      `gorm:"" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m *Mdl) MarshalBinary() ([]byte, error) {
	return sonic.Marshal(m)
}

func (m *Mdl) UnmarshalBinary(data []byte) error {
	return sonic.Unmarshal(data, &m)
}

func (m *Mdl) ToString() string {
	return fmt.Sprintf("%+v", m)
}
