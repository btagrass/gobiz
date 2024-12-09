package utl

import (
	"errors"
	"log/slog"
	"os/exec"
	"time"
)

func Command(cmds ...string) (string, error) {
	var output []byte
	var err error
	for _, c := range cmds {
		cs := Split(c, ' ')
		output, err = exec.Command(cs[0], cs[1:]...).CombinedOutput()
		if err != nil {
			return "", errors.New(string(output))
		}
	}
	return string(output), err
}

func ElapsedTime(method string, begin time.Time) {
	slog.Info("ElapsedTime", method, time.Since(begin).Milliseconds())
}
