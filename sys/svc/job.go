package svc

import (
	"time"

	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/robfig/cron/v3"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
)

type IJob interface {
	cron.Job
	GetName() string
	GetDesc() string
	GetCron() string
	GetArg() string
	GetArgDesc() string
	SetArg(arg string)
}

type JobSvc struct {
	*svc.DataSvc[mdl.Job]
	cron     *cron.Cron
	interval time.Duration
	jobs     map[string]IJob
}

func NewJobSvc(i *do.Injector) (*JobSvc, error) {
	s := &JobSvc{
		DataSvc:  svc.NewDataSvc[mdl.Job]("sys:jobs"),
		cron:     cron.New(cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger))),
		interval: viper.GetDuration("job.interval"),
		jobs:     make(map[string]IJob),
	}
	if s.interval == 0 {
		s.interval = 3 * time.Second
	}
	go func() {
		s.cron.Start()
		defer s.cron.Stop()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for range ticker.C {
			jobs, _, err := s.List()
			if err != nil {
				logrus.Error(err)
			}
			for _, j := range jobs {
				key := s.GetFullKey(j.Id)
				if j.State == 0 {
					if j.Instance == 0 {
						continue
					}
					s.cron.Remove(cron.EntryID(j.Instance))
					err = s.Update(map[string]any{
						"instance": 0,
					}, j.Id)
					if err != nil {
						logrus.Error(err)
						continue
					}
					s.Local.Delete(key)
				} else {
					v, ok := s.Local.Get(key)
					if ok {
						job := v.(mdl.Job)
						if job.Cron == j.Cron && job.Arg == j.Arg {
							entry := s.cron.Entry(cron.EntryID(j.Instance))
							err = s.Update(map[string]any{
								"updated_at": entry.Prev,
							}, j.Id)
							if err != nil {
								logrus.Error(err)
							}
							continue
						}
					}
					if j.Instance > 0 {
						s.cron.Remove(cron.EntryID(j.Instance))
					}
					job, ok := s.jobs[j.Name]
					if ok {
						job.SetArg(j.Arg)
					}
					instance, err := s.cron.AddJob(j.Cron, job)
					if err != nil {
						logrus.Error(err)
						continue
					}
					err = s.Update(map[string]any{
						"instance": instance,
					}, j.Id)
					if err != nil {
						logrus.Error(err)
						continue
					}
					s.Local.SetDefault(key, j)
				}
			}
		}
	}()
	return s, nil
}

func (s *JobSvc) AddJobs(jobs ...IJob) error {
	for _, j := range jobs {
		s.jobs[j.GetName()] = j
		err := s.Save(mdl.Job{
			Name:    j.GetName(),
			Desc:    j.GetDesc(),
			Cron:    j.GetCron(),
			Arg:     j.GetArg(),
			ArgDesc: j.GetArgDesc(),
		}, clause.OnConflict{
			Columns: []clause.Column{{
				Name: "name",
			}},
			DoUpdates: clause.AssignmentColumns([]string{
				"desc",
				"arg_desc",
				"updated_at",
			}),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *JobSvc) StartJob(id int64) error {
	return s.Update(map[string]any{
		"state": 1,
	}, id)
}

func (s *JobSvc) StopJob(id int64) error {
	return s.Update(map[string]any{
		"state": 0,
	}, id)
}
