package svc

import (
	"log/slog"
	"time"

	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/robfig/cron/v3"
	"github.com/samber/do"
	"github.com/spf13/viper"
	"gorm.io/gorm/clause"
)

type IJob interface {
	cron.Job
	Name() string
	Desc() string
	Cron() string
	Arg() string
	ArgDesc() string
	Stopping() bool
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
			jobs, _, _ := s.List()
			for _, j := range jobs {
				job, ok := s.jobs[j.Name]
				if !ok {
					continue
				}
				key := s.GetFullKey(j.Id)
				if j.State == 0 {
					if j.Instance == 0 {
						continue
					}
					s.cron.Remove(cron.EntryID(j.Instance))
					err := s.Update(map[string]any{
						"instance": 0,
					}, j.Id)
					if err != nil {
						slog.Error(err.Error())
					}
					s.Local.Delete(key)
				} else {
					if job.Stopping() {
						delete(s.jobs, job.Name())
						if j.Instance > 0 {
							s.cron.Remove(cron.EntryID(j.Instance))
						}
						s.Local.Delete(key)
						continue
					}
					v, ok := s.Local.Get(key)
					if ok {
						cj := v.(mdl.Job)
						if j.Cron == cj.Cron && j.Arg == cj.Arg {
							entry := s.cron.Entry(cron.EntryID(j.Instance))
							err := s.Update(map[string]any{
								"updated_at": entry.Prev,
							}, j.Id)
							if err != nil {
								slog.Error(err.Error())
							}
							continue
						}
					}
					if j.Instance > 0 {
						s.cron.Remove(cron.EntryID(j.Instance))
					}
					job.SetArg(j.Arg)
					instance, err := s.cron.AddJob(j.Cron, job)
					if err != nil {
						slog.Error(err.Error())
						continue
					}
					err = s.Update(map[string]any{
						"instance": instance,
					}, j.Id)
					if err != nil {
						slog.Error(err.Error())
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
		s.jobs[j.Name()] = j
		err := s.Save(mdl.Job{
			Name:    j.Name(),
			Desc:    j.Desc(),
			Cron:    j.Cron(),
			Arg:     j.Arg(),
			ArgDesc: j.ArgDesc(),
			State:   1,
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
