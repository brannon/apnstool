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

func NewSendBackgroundCommand() *cobra.Command {
	client := apns.NewClient()
	op := operation.NewSendBackground(client)

	var verbose bool

	cobraCmd := &cobra.Command{
		Use:   "background",
		Short: "Send simple background notification through APNs",
		RunE: func(c *cobra.Command, args []string) error {
			io := cmdio.NewCmdIO(c.OutOrStdout())

			if verbose {
				client.EnableLogging(io.Stdout())
			}

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
	flags.BoolVarP(&verbose, VerboseFlag, VerboseShortFlag, VerboseDefault, VerboseDesc)
	flags.StringVarP(&op.DataString, DataStringFlag, DataStringShortFlag, DataStringDefault, DataStringDesc)

	_ = cobraCmd.MarkFlagRequired(AppIdFlag)
	_ = cobraCmd.MarkFlagRequired(DeviceTokenFlag)

	return cobraCmd
}
