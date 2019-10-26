// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package auth

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "APNs authentication commands",
		Args:  cobra.NoArgs,
	}

	authCmd.AddCommand(NewAuthGenerateTokenCommand())

	return authCmd
}
