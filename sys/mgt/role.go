package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
)

func ListRoles(c *gin.Context) {
	roles, count, err := s.Use[*svc.RoleSvc]().List(r.Q(c))
	r.J(c, roles, count, err)
}

func GetRole(c *gin.Context) {
	role, err := s.Use[*svc.RoleSvc]().Get(c.Param("id"))
	r.J(c, role, err)
}

func RemoveRoles(c *gin.Context) {
	err := s.Use[*svc.RoleSvc]().Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

func SaveRole(c *gin.Context) {
	var role mdl.Role
	err := c.ShouldBind(&role)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.RoleSvc]().Save(role)
	r.J(c, role.Id, err)
}

func ListRoleResources(c *gin.Context) {
	resources, err := s.Use[*svc.RoleSvc]().ListRoleResources(c.Param("id"))
	r.J(c, resources, err)
}

func SaveRoleResources(c *gin.Context) {
	var resources []mdl.Resource
	err := c.ShouldBind(&resources)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.RoleSvc]().SaveRoleResources(c.Param("id"), resources)
	r.J(c, true, err)
}
