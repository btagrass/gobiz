package app

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/btagrass/gobiz/utl"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
)

var (
	Dir      string
	DataDir  string
	LogFile  io.Writer
	LogLevel slog.Leveler
)

func init() {
	// Config
	var err error
	Dir, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	confPath := filepath.Join(Dir, "conf/app.*")
	confFiles, err := filepath.Glob(confPath)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	for _, f := range confFiles {
		v := viper.New()
		v.SetConfigFile(f)
		err = v.ReadInConfig()
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		viper.MergeConfigMap(v.AllSettings())
	}
	DataDir = viper.GetString("data.dir")
	if DataDir == "" {
		DataDir = filepath.Join(Dir, "data")
	} else if filepath.IsLocal(DataDir) {
		DataDir = filepath.Join(Dir, DataDir)
	}
	err = utl.MakeDir(DataDir)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	// Log
	logPath := filepath.Join(Dir, "logs/%Y%m%d.log")
	logCount := viper.GetUint("log.count")
	if logCount == 0 {
		logCount = 7
	}
	LogFile, err = rotatelogs.New(
		logPath,
		rotatelogs.WithMaxAge(-1),
		rotatelogs.WithRotationCount(logCount),
	)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	switch strings.ToLower(viper.GetString("log.level")) {
	case "info":
		LogLevel = slog.LevelInfo
	case "warn":
		LogLevel = slog.LevelWarn
	case "error":
		LogLevel = slog.LevelError
	default:
		LogLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(io.MultiWriter(os.Stdout, LogFile), &slog.HandlerOptions{
		AddSource: true,
		Level:     LogLevel,
	})))
}
