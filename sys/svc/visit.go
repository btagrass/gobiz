package svc

import (
	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/samber/do"
)

type VisitSvc struct {
	*svc.DataSvc[mdl.Visit]
}

func NewVisitSvc(i *do.Injector) (*VisitSvc, error) {
	return &VisitSvc{
		DataSvc: svc.NewDataSvc[mdl.Visit]("sys:visits"),
	}, nil
}
