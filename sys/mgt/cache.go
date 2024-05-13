package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
)

func ListCaches(c *gin.Context) {
	caches, count, err := s.Use[*svc.CacheSvc]().ListCaches(c.Query("type"), c.Query("keyword"))
	r.J(c, caches, count, err)
}

func RemoveCaches(c *gin.Context) {
	err := s.Use[*svc.CacheSvc]().RemoveCaches(c.Query("type"), utl.Split(c.Param("ids"), ',')...)
	r.J(c, true, err)
}
