package svc

import (
	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type DeptSvc struct {
	*svc.DataSvc[mdl.Dept]
}

func NewDeptSvc(i *do.Injector) (*DeptSvc, error) {
	return &DeptSvc{
		DataSvc: svc.NewDataSvc[mdl.Dept]("sys:depts"),
	}, nil
}

func (s *DeptSvc) ListDepts(conds map[string]any) ([]mdl.Dept, int64, error) {
	var depts []mdl.Dept
	var count int64
	db := s.
		Make(conds).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Where("parent_id = 0").
		Order("sequence").
		Find(&depts)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	err := db.Error
	if err != nil {
		return depts, count, err
	}
	return depts, count, nil
}
