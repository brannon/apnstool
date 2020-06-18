// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package send

import (
	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/cmdio"
	"github.com/brannon/apnstool/operation"
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

func NewSendAlertCommand() *cobra.Command {
	client := apns.NewClient()
	op := operation.NewSendAlert(client)

	cobraCmd := &cobra.Command{
		Use:   "alert",
		Short: "Send simple alert notification through APNs",
		RunE: func(c *cobra.Command, args []string) error {
			io := cmdio.NewCmdIO(c.OutOrStdout())

			result, err := op.Exec()
			if err != nil {
				return err
			}

			io.Out("Notification sent successfully\n")
			io.Outf("APNS-ID: %s\n", result.ApnsId)

			return nil
		},
	}

	flags := cobraCmd.Flags()
	BindSendOperationCommonFlags(flags, &op.SendOperation)
	flags.StringVar(&op.AlertText, AlertTextFlag, AlertTextDefault, AlertTextDesc)
	flags.IntVar(&op.BadgeCount, BadgeCountFlag, BadgeCountDefault, BadgeCountDesc)
	flags.StringVar(&op.SoundName, SoundNameFlag, SoundNameDefault, SoundNameDesc)

	_ = cobraCmd.MarkFlagRequired(AppIdFlag)
	_ = cobraCmd.MarkFlagRequired(DeviceTokenFlag)

	return cobraCmd
}
