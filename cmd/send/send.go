// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package send

import (
	"github.com/brannon/apnstool/cmd/auth"
	"github.com/brannon/apnstool/operation"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	AppIdFlag    = "app-id"
	AppIdDefault = ""
	AppIdDesc    = "app bundle ID"

	DataStringFlag      = "data"
	DataStringShortFlag = "d"
	DataStringDefault   = ""
	DataStringDesc      = "JSON formatted notification content"

	DeviceTokenFlag    = "device-token"
	DeviceTokenDefault = ""
	DeviceTokenDesc    = "APNs device token"

	SandboxFlag    = "sandbox"
	SandboxDefault = false
	SandboxDesc    = "use APNS sandbox endpoint"

	VerboseFlag      = "verbose"
	VerboseShortFlag = "v"
	VerboseDefault   = false
	VerboseDesc      = "enable verbose logging"
)

func BindSendOperationCommonFlags(flags *pflag.FlagSet, op *operation.SendOperation) {
	auth.BindTokenAuthFlags(flags, &op.TokenAuth)
	auth.BindCertificateAuthFlags(flags, &op.CertificateAuth)
	flags.StringVar(&op.AppId, AppIdFlag, AppIdDefault, AppIdDesc)
	flags.StringVar(&op.DeviceToken, DeviceTokenFlag, DeviceTokenDefault, DeviceTokenDesc)
	flags.BoolVar(&op.Sandbox, SandboxFlag, SandboxDefault, SandboxDesc)
}

func GetCommand() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:   "send",
		Short: "APNs send notification commands",
	}

	sendCmd.AddCommand(NewSendAlertCommand())
	sendCmd.AddCommand(NewSendBackgroundCommand())
	sendCmd.AddCommand(NewSendRawCommand())

	return sendCmd
}
