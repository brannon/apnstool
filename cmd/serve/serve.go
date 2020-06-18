// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"net/http"

	"github.com/spf13/cobra"
)

type ServeCmd struct {
}

func GetCommand() *cobra.Command {
	cmd := &ServeCmd{}

	http.HandleFunc("/send/alert", cmd.handlePostSendAlert)
	http.HandleFunc("/", cmd.handleGetIndex)

	cobraCmd := &cobra.Command{
		Use:   "serve",
		Short: "Run tool as web server",
		RunE: func(c *cobra.Command, args []string) error {
			return cmd.Run()
		},
	}

	return cobraCmd
}

func (cmd *ServeCmd) Run() error {
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		return err
	}
	return nil
}
