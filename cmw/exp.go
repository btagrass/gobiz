package cmw

import (
	"fmt"
	"time"

	"github.com/btagrass/gobiz/r"
	"github.com/gin-gonic/gin"
)

func Exp(date string) gin.HandlerFunc {
	return func(c *gin.Context) {
		expirationDate, _ := time.Parse(time.DateOnly, date)
		if time.Now().After(expirationDate) {
			r.J(c, fmt.Errorf("license has expired"))
		}
		c.Next()
	}
}
