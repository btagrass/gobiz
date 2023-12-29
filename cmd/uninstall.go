package cmd

import (
	"fmt"

	"github.com/btagrass/gobiz/utl"
	"github.com/spf13/cobra"
)

var (
	Uninstall = &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall",
		Run: func(c *cobra.Command, args []string) {
			name := cmd.Use
			_, err := utl.Command(fmt.Sprintf("systemctl stop %s", name), fmt.Sprintf("systemctl disable %s", name))
			if err != nil {
				fmt.Print(err)
				return
			}
			err = utl.Remove(fmt.Sprintf("/etc/systemd/system/%s.service", name))
			if err != nil {
				fmt.Print(err)
			}
		},
	}
)
