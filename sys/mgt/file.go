package mgt

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
)

func SaveFile(c *gin.Context) {
	header, err := c.FormFile("file")
	if err != nil {
		c.Abort()
	}
	fileDir := filepath.Join("data/files", c.Param("dir"))
	err = utl.MakeDir(fileDir)
	if err != nil {
		c.Abort()
	}
	fileName := fmt.Sprintf("%s/%s%s", fileDir, utl.TimeId(), filepath.Ext(header.Filename))
	err = c.SaveUploadedFile(header, fileName)
	if err != nil {
		c.Abort()
	}

	c.String(http.StatusOK, fileName)
}
