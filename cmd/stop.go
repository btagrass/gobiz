package cmd

import (
	"fmt"

	"github.com/btagrass/gobiz/utl"
	"github.com/spf13/cobra"
)

var (
	Stop = &cobra.Command{
		Use:   "stop",
		Short: "Stop",
		Run: func(c *cobra.Command, args []string) {
			name := cmd.Use
			_, err := utl.Command(fmt.Sprintf("systemctl stop %s", name))
			if err != nil {
				fmt.Print(err)
			}
		},
	}
)
