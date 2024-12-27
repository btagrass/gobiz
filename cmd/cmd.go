package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{}
)

func Execute(name, desc string, cmds ...*cobra.Command) {
	cmd.Use = name
	cmd.Short = desc
	for _, c := range cmds {
		cmd.AddCommand(c)
	}
	err := cmd.Execute()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
