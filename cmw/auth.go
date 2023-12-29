package cmw

import (
	"fmt"

	"github.com/btagrass/gobiz/r"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/cast"
)

func Auth(perm *casbin.SyncedEnforcer, signedKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			r.J(c, fmt.Errorf("token is invalid"))
			return
		}
		token, err := jwt.Parse(authorization, func(token *jwt.Token) (any, error) {
			return signedKey, nil
		})
		if err != nil || !token.Valid {
			r.J(c, fmt.Errorf("token is invalid"))
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			r.J(c, fmt.Errorf("token is invalid"))
			return
		}
		userId, ok := claims["userId"]
		if !ok {
			r.J(c, fmt.Errorf("token is invalid"))
			return
		}
		ok, err = perm.Enforce(cast.ToString(userId), c.Request.URL.Path, c.Request.Method)
		if err != nil || !ok {
			r.J(c, fmt.Errorf("no permission to access %s", c.Request.URL.Path))
			return
		}
		c.Set("userId", userId)
		c.Set("userName", claims["userName"])
		c.Next()
	}
}
