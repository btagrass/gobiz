package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
)

func ListDepts(c *gin.Context) {
	depts, count, err := s.Use[*svc.DeptSvc]().ListDepts(r.Q(c))
	r.J(c, depts, count, err)
}

func GetDept(c *gin.Context) {
	dept, err := s.Use[*svc.DeptSvc]().Get(c.Param("id"))
	r.J(c, dept, err)
}

func RemoveDepts(c *gin.Context) {
	err := s.Use[*svc.DeptSvc]().Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

func SaveDept(c *gin.Context) {
	var dept mdl.Dept
	err := c.ShouldBind(&dept)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.DeptSvc]().Save(dept)
	r.J(c, dept.Id, err)
}
