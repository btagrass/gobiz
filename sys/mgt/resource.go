package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func ListResources(c *gin.Context) {
	resources, count, err := s.Use[*svc.ResourceSvc]().ListResources(r.Q(c))
	r.J(c, resources, count, err)
}

func GetResource(c *gin.Context) {
	resource, err := s.Use[*svc.ResourceSvc]().Get(c.Param("id"))
	r.J(c, resource, err)
}

func RemoveResources(c *gin.Context) {
	err := s.Use[*svc.ResourceSvc]().Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

func SaveResource(c *gin.Context) {
	var resource mdl.Resource
	err := c.ShouldBind(&resource)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.ResourceSvc]().Save(resource)
	r.J(c, resource.Id, err)
}

func ListMenus(c *gin.Context) {
	menus, err := s.Use[*svc.ResourceSvc]().ListMenus(cast.ToString(c.GetFloat64("userId")))
	r.J(c, menus, err)
}
