package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/gin-gonic/gin"
)

func GetServer(c *gin.Context) {
	server := s.Use[*svc.ServerSvc]().GetServer()
	r.J(c, server)
}
