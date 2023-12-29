package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
)

func ListVisits(c *gin.Context) {
	visits, count, err := s.Use[*svc.VisitSvc]().List(r.Q(c))
	r.J(c, visits, count, err)
}

func RemoveVisits(c *gin.Context) {
	err := s.Use[*svc.VisitSvc]().Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

func ClearVisits(c *gin.Context) {
	err := s.Use[*svc.VisitSvc]().Purge()
	r.J(c, true, err)
}
