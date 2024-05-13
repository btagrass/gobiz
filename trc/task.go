package trc

import (
	"time"
)

type ITask interface {
	GetCode() string
	GetBeginTime() time.Time
	GetEndTime() time.Time
	GetCount() int
}
