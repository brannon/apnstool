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
	AlertTextFlag    = "alert-text"
	AlertTextDefault = ""
	AlertTextDesc    = "alert text"

	BadgeCountFlag    = "badge-count"
	BadgeCountDefault = 0
	BadgeCountDesc    = "badge count"

	SoundNameFlag    = "sound-name"
	SoundNameDefault = ""
	SoundNameDesc    = "sound name"
)

type SendAlertCmd struct {
	SendCmd

	AlertText  string
	BadgeCount int
	SoundName  string
}

func NewSendAlertCommand() *cobra.Command {
	cmd := &SendAlertCmd{}

	cobraCmd := &cobra.Command{
		Use:   "alert",
		Short: "Send simple alert notification through APNs",
		RunE: func(c *cobra.Command, args []string) error {
			cmd.Client = apns.NewClient()
			cmd.IO = cmdio.NewCmdIO(c.OutOrStdout())

			return cmd.Run()
		},
	}

	flags := cobraCmd.Flags()
	BindSendCommonFlags(flags, &cmd.SendCmd)
	flags.StringVar(&cmd.AlertText, AlertTextFlag, AlertTextDefault, AlertTextDesc)
	flags.IntVar(&cmd.BadgeCount, BadgeCountFlag, BadgeCountDefault, BadgeCountDesc)
	flags.StringVar(&cmd.SoundName, SoundNameFlag, SoundNameDefault, SoundNameDesc)

	return cobraCmd
}

func (cmd *SendAlertCmd) Run() error {
	notificationBuilder := apns.NewNotificationBuilder(cmd.AppId)

	if cmd.AlertText != "" {
		notificationBuilder.SetAlertText(cmd.AlertText)
	}

	if cmd.BadgeCount != -1 {
		notificationBuilder.SetBadgeCount(cmd.BadgeCount)
	}

	if cmd.SoundName != "" {
		notificationBuilder.SetSoundName(cmd.SoundName)
	}

	headers, content, err := notificationBuilder.Build()
	if err != nil {
		return err
	}

	return cmd.sendNotification(headers, content)
}
