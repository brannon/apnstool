// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package send

import (
	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/cmdio"
	"github.com/spf13/cobra"
)

type SendBackgroundCmd struct {
	SendCmd

	DataString string
}

func NewSendBackgroundCommand() *cobra.Command {
	cmd := &SendBackgroundCmd{}

	cobraCmd := &cobra.Command{
		Use:   "background",
		Short: "Send simple background notification through APNs",
		RunE: func(c *cobra.Command, args []string) error {
			cmd.Client = apns.NewClient()
			cmd.IO = cmdio.NewCmdIO(c.OutOrStdout())

			return cmd.Run()
		},
	}

	flags := cobraCmd.Flags()
	BindSendCommonFlags(flags, &cmd.SendCmd)
	flags.StringVarP(&cmd.DataString, DataStringFlag, DataStringShortFlag, DataStringDefault, DataStringDesc)

	_ = cobraCmd.MarkFlagRequired(AppIdFlag)
	_ = cobraCmd.MarkFlagRequired(DeviceTokenFlag)

	return cobraCmd
}

func (cmd *SendBackgroundCmd) Run() error {
	notificationBuilder := apns.NewNotificationBuilder(cmd.AppId)
	notificationBuilder.SetContentAvailable(true)

	if cmd.DataString != "" {
		data, err := parseDataString(cmd.DataString)
		if err != nil {
			return err
		}

		notificationBuilder.Merge(data)
	}

	headers, content, err := notificationBuilder.Build()
	if err != nil {
		return err
	}

	return cmd.sendNotification(headers, content)
}
