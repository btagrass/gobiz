package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/btagrass/gobiz/utl"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	users, count, err := s.Use[*svc.UserSvc]().ListUsers(r.Q(c))
	r.J(c, users, count, err)
}

func GetUser(c *gin.Context) {
	user, err := s.Use[*svc.UserSvc]().Get(c.Param("id"))
	r.J(c, user, err)
}

func RemoveUsers(c *gin.Context) {
	err := s.Use[*svc.UserSvc]().RemoveUsers(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

func SaveUser(c *gin.Context) {
	var user mdl.User
	err := c.ShouldBind(&user)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.UserSvc]().SaveUser(user)
	r.J(c, user.Id, err)
}

func ListUserRoles(c *gin.Context) {
	roles, err := s.Use[*svc.UserSvc]().ListUserRoles(c.Param("id"))
	r.J(c, roles, err)
}

func SaveUserRoles(c *gin.Context) {
	var roles []int64
	err := c.ShouldBind(&roles)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.UserSvc]().SaveUserRoles(c.Param("id"), roles)
	r.J(c, true, err)
}

func Login(c *gin.Context) {
	var p struct {
		UserName string `form:"userName" json:"userName" binding:"required"` // 用户名
		Password string `form:"password" json:"password" binding:"required"` // 密码
	}
	err := c.ShouldBind(&p)
	if err != nil {
		r.J(c, err)
		return
	}
	user, err := s.Use[*svc.UserSvc]().Login(p.UserName, p.Password)
	r.J(c, user, err)
}
