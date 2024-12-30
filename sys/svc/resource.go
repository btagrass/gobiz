package svc

import (
	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type ResourceSvc struct {
	*svc.DataSvc[mdl.Resource]
	userSvc *UserSvc
}

func NewResourceSvc(i *do.Injector) (*ResourceSvc, error) {
	return &ResourceSvc{
		DataSvc: svc.NewDataSvc[mdl.Resource]("sys:resources"),
		userSvc: svc.Use[*UserSvc](),
	}, nil
}

func (s *ResourceSvc) ListMenus(userId string) ([]mdl.Resource, error) {
	var resources []mdl.Resource
	err := s.Make().
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Where("type = 1").Order("sequence")
		}).
		Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("Children.Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Where("parent_id = 0 and type = 1").
		Order("sequence").
		Find(&resources).Error
	if err != nil {
		return resources, err
	}
	if userId == "300000000000001" {
		return resources, nil
	}
	var resourceIds []string
	roles, _ := s.userSvc.Perm.GetRolesForUser(userId)
	for _, r := range roles {
		permissions, err := s.userSvc.Perm.GetPermissionsForUser(r)
		if err != nil {
			return nil, err
		}
		for _, p := range permissions {
			resourceIds = append(resourceIds, p[3])
		}
	}
	for i := 0; i < len(resources); i++ {
		r := &resources[i]
		if lo.Contains(resourceIds, cast.ToString(r.Id)) {
			continue
		}
		for i1 := 0; i1 < len(r.Children); i1++ {
			r1 := &r.Children[i1]
			if lo.Contains(resourceIds, cast.ToString(r1.Id)) {
				continue
			}
			for i2 := 0; i2 < len(r1.Children); i2++ {
				r2 := &r1.Children[i2]
				if lo.Contains(resourceIds, cast.ToString(r2.Id)) {
					continue
				}
				for i3 := 0; i3 < len(r2.Children); i3++ {
					r3 := &r2.Children[i2]
					if !lo.Contains(resourceIds, cast.ToString(r3.Id)) {
						r2.Children = append(r2.Children[:i3], r2.Children[i3+1:]...)
						i3--
					}
				}
				if len(r2.Children) == 0 {
					r1.Children = append(r1.Children[:i2], r1.Children[i2+1:]...)
					i2--
				}
			}
			if len(r1.Children) == 0 {
				r.Children = append(r.Children[:i1], r.Children[i1+1:]...)
				i1--
			}
		}
		if len(r.Children) == 0 {
			resources = append(resources[:i], resources[i+1:]...)
			i--
		}
	}
	return resources, nil
}

func (s *ResourceSvc) ListResources(conds map[string]any) ([]mdl.Resource, int64, error) {
	var resources []mdl.Resource
	var count int64
	db := s.
		Make(conds).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Preload("Children.Children.Children", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence")
		}).
		Where("parent_id = 0").
		Order("sequence").
		Find(&resources)
	_, ok := db.Statement.Clauses["LIMIT"]
	if ok {
		db = db.Limit(-1).Offset(-1).Count(&count)
	}
	err := db.Error
	if err != nil {
		return resources, count, err
	}
	return resources, count, nil
}
