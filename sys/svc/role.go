package svc

import (
	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/samber/do"
	"github.com/spf13/cast"
)

type RoleSvc struct {
	*svc.DataSvc[mdl.Role]
	userSvc *UserSvc
}

func NewRoleSvc(i *do.Injector) (*RoleSvc, error) {
	return &RoleSvc{
		DataSvc: svc.NewDataSvc[mdl.Role]("sys:roles"),
		userSvc: svc.Use[*UserSvc](),
	}, nil
}

func (s *RoleSvc) ListRoleResources(id string) ([]int64, error) {
	resources := make([]int64, 0)
	permissions, err := s.userSvc.Perm.GetPermissionsForUser(id)
	if err != nil {
		return nil, err
	}
	for _, p := range permissions {
		resources = append(resources, cast.ToInt64(p[3]))
	}
	return resources, nil
}

func (s *RoleSvc) SaveRoleResources(id string, resources []mdl.Resource) error {
	_, err := s.userSvc.Perm.DeletePermissionsForUser(id)
	if err != nil {
		return err
	}
	if len(resources) == 0 {
		return nil
	}
	var rs [][]string
	for _, r := range resources {
		rs = append(rs, []string{
			r.Uri,
			r.Act,
			cast.ToString(r.Id),
		})
	}
	_, err = s.userSvc.Perm.AddPermissionsForUser(id, rs...)
	return err
}
