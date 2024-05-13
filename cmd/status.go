package cmd

import (
	"fmt"

	"github.com/btagrass/gobiz/utl"
	"github.com/spf13/cobra"
)

var (
	Status = &cobra.Command{
		Use:   "status",
		Short: "Status",
		Run: func(c *cobra.Command, args []string) {
			name := cmd.Use
			_, err := utl.Command(fmt.Sprintf("systemctl status %s", name))
			if err != nil {
				fmt.Print(err)
			}
		},
	}
)
