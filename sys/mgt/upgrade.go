package mgt

import (
	"net/http"
	"path/filepath"

	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/gin-gonic/gin"
)

func Upgrade(c *gin.Context) {
	filePath, fileVer, err := s.Use[*svc.UpgradeSvc]().Upgrade(c.Param("ver"))
	if err != nil || filePath == "" {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Header("ver", fileVer)
		c.FileAttachment(filePath, filepath.Base(filePath))
	}
}
