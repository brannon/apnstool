// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"fmt"
	"net/http"

	"github.com/brannon/apnstool/cmdio"
	"github.com/spf13/cobra"
)

const (
	LocalOnlyFlag    = "local"
	LocalOnlyDefault = false
	LocalOnlyDesc    = "Restrict incoming connections to localhost. Useful for avoiding firewall prompts when testing."

	PortFlag    = "port"
	PortDefault = 8080
	PortDesc    = "TCP port on which to listen"
)

type ServeCmd struct {
	Port      int
	LocalOnly bool

	IO cmdio.CmdIO
}

func GetCommand() *cobra.Command {
	cmd := &ServeCmd{}

	http.HandleFunc("/send/alert", cmd.handleSendAlert)
	http.HandleFunc("/send/background", cmd.handleSendBackground)
	http.HandleFunc("/send/raw", cmd.handleSendRaw)
	http.HandleFunc("/", cmd.handleIndex)

	cobraCmd := &cobra.Command{
		Use:   "serve",
		Short: "Run tool as web server",
		RunE: func(c *cobra.Command, args []string) error {
			cmd.IO = cmdio.NewCmdIO(c.OutOrStdout())

			return cmd.Run()
		},
	}

	flags := cobraCmd.Flags()
	flags.BoolVar(&cmd.LocalOnly, LocalOnlyFlag, LocalOnlyDefault, LocalOnlyDesc)
	flags.IntVar(&cmd.Port, PortFlag, PortDefault, PortDesc)

	return cobraCmd
}

func (cmd *ServeCmd) Run() error {
	hostname := ""
	if cmd.LocalOnly {
		hostname = "localhost"
	}

	addr := fmt.Sprintf("%s:%d", hostname, cmd.Port)

	cmd.IO.Outf("HTTP server listening on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}
	return nil
}
