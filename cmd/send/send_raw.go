// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package send

import (
	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/cmdio"
	"github.com/spf13/cobra"
)

const (
	PriorityFlag    = "priority"
	PriorityDefault = ""
	PriorityDesc    = "value for 'apns-priority' header"

	PushTypeFlag    = "push-type"
	PushTypeDefault = ""
	PushTypeDesc    = "value for 'apns-push-type' header"
)

type SendRawCmd struct {
	SendCmd

	DataString string
	Priority   string
	PushType   string
}

func NewSendRawCommand() *cobra.Command {
	cmd := &SendRawCmd{}

	cobraCmd := &cobra.Command{
		Use:   "raw",
		Short: "Send raw notification through APNs",
		RunE: func(c *cobra.Command, args []string) error {
			cmd.Client = apns.NewClient()
			cmd.IO = cmdio.NewCmdIO(c.OutOrStdout())

			return cmd.Run()
		},
	}

	flags := cobraCmd.Flags()
	BindSendCommonFlags(flags, &cmd.SendCmd)
	flags.StringVarP(&cmd.DataString, DataStringFlag, DataStringShortFlag, DataStringDefault, DataStringDesc)
	flags.StringVar(&cmd.Priority, PriorityFlag, PriorityDefault, PriorityDesc)
	flags.StringVar(&cmd.PushType, PushTypeFlag, PushTypeDefault, PushTypeDesc)

	return cobraCmd
}

func (cmd *SendRawCmd) Run() error {
	headers := make(apns.Headers)

	headers["apns-topic"] = cmd.AppId

	if cmd.Priority != "" {
		headers["apns-priority"] = cmd.Priority
	}

	if cmd.PushType != "" {
		headers["apns-push-type"] = cmd.PushType
	}

	return cmd.sendNotification(headers, []byte(cmd.DataString))
}
