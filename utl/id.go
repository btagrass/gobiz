package utl

import (
	"time"

	"github.com/yitter/idgenerator-go/idgen"
)

func init() {
	workerId := uint16(time.Now().Unix() % 64)
	options := idgen.NewIdGeneratorOptions(workerId)
	idgen.SetIdGenerator(options)
}

func IntId() int64 {
	return idgen.NextId()
}

func TimeId() string {
	return time.Now().Format("20060102150405.999999999")
}
