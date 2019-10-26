// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package send

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/cmd/auth"
	"github.com/brannon/apnstool/cmdio"

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

	VerboseFlag      = "verbose"
	VerboseShortFlag = "v"
	VerboseDefault   = false
	VerboseDesc      = "enable verbose logging"
)

type SendCmd struct {
	AppId       string
	DeviceToken string
	TokenAuth   auth.TokenAuth
	Verbose     bool

	Client apns.Client
	IO     cmdio.CmdIO
}

func BindSendCommonFlags(flags *pflag.FlagSet, cmd *SendCmd) {
	auth.BindTokenAuthFlags(flags, &cmd.TokenAuth)
	flags.StringVar(&cmd.AppId, AppIdFlag, AppIdDefault, AppIdDesc)
	flags.StringVar(&cmd.DeviceToken, DeviceTokenFlag, DeviceTokenDefault, DeviceTokenDesc)
	flags.BoolVarP(&cmd.Verbose, VerboseFlag, VerboseShortFlag, VerboseDefault, VerboseDesc)
}

func (cmd *SendCmd) sendNotification(
	headers apns.Headers,
	content []byte,
) error {
	if cmd.Verbose {
		cmd.Client.EnableLogging(cmd.IO.Stdout())
	}

	if cmd.useTokenAuth() {
		token, err := apns.GenerateJWTFromKeyFile(
			cmd.TokenAuth.KeyFile,
			cmd.TokenAuth.KeyId,
			cmd.TokenAuth.TeamId,
			time.Now(),
			cmd.TokenAuth.ExpiresAfter,
		)
		if err != nil {
			return err
		}

		cmd.Client.ConfigureTokenAuth(token)
	}

	result, err := cmd.Client.Send(cmd.DeviceToken, headers, content)
	if err != nil {
		return err
	}

	if result.Success() {
		cmd.IO.Out("Notification sent successfully\n")
		cmd.IO.Outf("APNS-ID: %s\n", result.Id())
	}

	return nil
}

func (cmd *SendCmd) useTokenAuth() bool {
	return cmd.TokenAuth.KeyFile != "" &&
		cmd.TokenAuth.KeyId != "" &&
		cmd.TokenAuth.TeamId != ""
}

func parseDataString(dataString string) (map[string]interface{}, error) {
	data := make(map[string]interface{})

	reader := bytes.NewBufferString(dataString)
	decoder := json.NewDecoder(reader)

	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
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
