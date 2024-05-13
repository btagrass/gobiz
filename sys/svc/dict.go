package svc

import (
	"github.com/btagrass/gobiz/svc"
	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/samber/do"
)

type DictSvc struct {
	*svc.DataSvc[mdl.Dict]
}

func NewDictSvc(i *do.Injector) (*DictSvc, error) {
	return &DictSvc{
		DataSvc: svc.NewDataSvc[mdl.Dict]("sys:dicts"),
	}, nil
}
