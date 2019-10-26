// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"time"

	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/cmdio"
	"github.com/spf13/cobra"
)

type AuthGenerateTokenCmd struct {
	TokenAuth TokenAuth

	io cmdio.CmdIO
}

func NewAuthGenerateTokenCommand() *cobra.Command {
	authGenerateToken := &AuthGenerateTokenCmd{}

	cobraCmd := &cobra.Command{
		Use:   "generate-token",
		Short: "Generate JWT token from .p8 key",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			authGenerateToken.io = cmdio.NewCmdIO(cmd.OutOrStdout())

			err := authGenerateToken.Prepare(args)
			if err != nil {
				return err
			}
			return authGenerateToken.Run()
		},
	}

	BindTokenAuthFlags(cobraCmd.Flags(), &authGenerateToken.TokenAuth)
	_ = cobraCmd.MarkFlagRequired(KeyFileFlag)
	_ = cobraCmd.MarkFlagRequired(KeyIdFlag)
	_ = cobraCmd.MarkFlagRequired(TeamIdFlag)

	return cobraCmd
}

func (cmd *AuthGenerateTokenCmd) Prepare(args []string) error {
	return nil
}

func (cmd *AuthGenerateTokenCmd) Run() error {
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

	cmd.io.Outf("%s\n", token)

	return nil
}
