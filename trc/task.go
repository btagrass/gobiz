package trc

import (
	"time"
)

type ITask interface {
	Code() string
	BeginTime() time.Time
	EndTime() time.Time
	Count() int
}
