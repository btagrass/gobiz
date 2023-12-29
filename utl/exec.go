package utl

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/sirupsen/logrus"
)

func Command(cmds ...string) (string, error) {
	var output []byte
	var err error
	for _, c := range cmds {
		cs := Split(c, ' ')
		output, err = exec.Command(cs[0], cs[1:]...).CombinedOutput()
		if err != nil {
			return "", fmt.Errorf(string(output))
		}
	}
	return string(output), err
}

func ElapsedTime(method string, begin time.Time) {
	elapsedTime := time.Since(begin)
	logrus.Infof("method: %s == %dms", method, elapsedTime.Milliseconds())
}
