package mgt

import (
	"github.com/btagrass/gobiz/cmw"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Mgt() *gin.Engine {
	e := gin.Default()
	// Cors
	e.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*", "Authorization"},
	}))
	// Debug
	pprof.Register(e)
	// Mgt
	m := e.Group("/mgt")
	{
		m.POST("/login", Login)
		m.GET("/upgrades/:ver", Upgrade)
	}
	// Sys
	s := m.Group("/sys").Use(cmw.Auth(s.Use[*svc.UserSvc]().Perm, s.Use[*svc.UserSvc]().SignedKey), cmw.Visit(s.Use[*svc.VisitSvc]()))
	{
		s.GET("/caches", ListCaches)
		s.DELETE("/caches/:ids", RemoveCaches)

		s.GET("/depts", ListDepts)
		s.GET("/depts/:id", GetDept)
		s.DELETE("/depts/:ids", RemoveDepts)
		s.POST("/depts", SaveDept)

		s.GET("/dicts", ListDicts)
		s.GET("/dicts/:id", GetDict)
		s.DELETE("/dicts/:ids", RemoveDicts)
		s.POST("/dicts", SaveDict)

		s.POST("/files/:dir", SaveFile)

		s.GET("/jobs", ListJobs)
		s.GET("/jobs/:id", GetJob)
		s.POST("/jobs", SaveJob)
		s.POST("/jobs/:id/start", StartJob)
		s.POST("/jobs/:id/stop", StopJob)

		s.GET("/resources", ListResources)
		s.GET("/resources/:id", func(c *gin.Context) {
			if c.Param("id") == "menu" {
				ListMenus(c)
			} else {
				GetResource(c)
			}
		})
		s.DELETE("/resources/:ids", RemoveResources)
		s.POST("/resources", SaveResource)

		s.GET("/roles", ListRoles)
		s.GET("/roles/:id", GetRole)
		s.DELETE("/roles/:ids", RemoveRoles)
		s.POST("/roles", SaveRole)
		s.GET("/roles/:id/resources", ListRoleResources)
		s.POST("/roles/:id/resources", SaveRoleResources)

		s.GET("/servers", GetServer)

		s.GET("/users", ListUsers)
		s.GET("/users/:id", GetUser)
		s.DELETE("/users/:ids", RemoveUsers)
		s.POST("/users", SaveUser)
		s.GET("/users/:id/roles", ListUserRoles)
		s.POST("/users/:id/roles", SaveUserRoles)

		s.GET("/visits", ListVisits)
		s.DELETE("/visits/:ids", func(c *gin.Context) {
			if c.Param("ids") == "all" {
				ClearVisits(c)
			} else {
				RemoveVisits(c)
			}
		})
	}
	return e
}
