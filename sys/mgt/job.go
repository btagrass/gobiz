package mgt

import (
	"github.com/btagrass/gobiz/r"
	s "github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func ListJobs(c *gin.Context) {
	jobs, count, err := s.Use[*svc.JobSvc]().List(r.Q(c))
	r.J(c, jobs, count, err)
}

func GetJob(c *gin.Context) {
	job, err := s.Use[*svc.JobSvc]().Get(c.Param("id"))
	r.J(c, job, err)
}

func SaveJob(c *gin.Context) {
	var job mdl.Job
	err := c.ShouldBind(&job)
	if err != nil {
		r.J(c, err)
		return
	}
	err = s.Use[*svc.JobSvc]().Save(job)
	r.J(c, job.Id, err)
}

func StartJob(c *gin.Context) {
	err := s.Use[*svc.JobSvc]().StartJob(cast.ToInt64(c.Param("id")))
	r.J(c, true, err)
}

func StopJob(c *gin.Context) {
	err := s.Use[*svc.JobSvc]().StopJob(cast.ToInt64(c.Param("id")))
	r.J(c, true, err)
}
