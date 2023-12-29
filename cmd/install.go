package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/btagrass/gobiz/app"
	"github.com/btagrass/gobiz/utl"
	"github.com/spf13/cobra"
)

var (
	Install = &cobra.Command{
		Use:   "install",
		Short: "Install",
		Run: func(c *cobra.Command, args []string) {
			name := cmd.Use
			err := os.WriteFile(fmt.Sprintf("/etc/systemd/system/%s.service", name), []byte(fmt.Sprintf(`
[Unit]
Description=%s
After=network.target

[Service]
Type=simple
WorkingDirectory=%s
ExecStart=%s run
Restart=always
RestartSec=30s

[Install]
WantedBy=multi-user.target
`, strings.ToUpper(name), app.Dir, filepath.Join(app.Dir, name))), os.ModePerm)
			if err != nil {
				fmt.Print(err)
				return
			}
			_, err = utl.Command(fmt.Sprintf("systemctl enable %s", name), "systemctl daemon-reload", fmt.Sprintf("systemctl restart %s", name))
			if err != nil {
				fmt.Print(err)
			}
		},
	}
)
