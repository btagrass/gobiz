package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
)

func ListDicts(c *gin.Context) {
	dicts, count, err := s.Use[*svc.DictSvc]().List(r.Q(c))
	r.J(c, dicts, count, err)
}

func GetDict(c *gin.Context) {
	dict, err := s.Use[*svc.DictSvc]().Get(c.Param("id"))
	r.J(c, dict, err)
}

func RemoveDicts(c *gin.Context) {
	err := s.Use[*svc.DictSvc]().Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

func SaveDict(c *gin.Context) {
	var dict mdl.Dict
	err := c.ShouldBind(&dict)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.DictSvc]().Save(dict)
	r.J(c, dict.Id, err)
}
