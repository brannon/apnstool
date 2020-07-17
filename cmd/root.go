// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/brannon/apnstool/build_version"
	"github.com/brannon/apnstool/cmd/auth"
	"github.com/brannon/apnstool/cmd/send"
	"github.com/brannon/apnstool/cmd/serve"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "apnstool",
	Short:        "APNSTool is a command-line tool for interacting with APNs",
	Version:      fmt.Sprintf("%s-%s", build_version.BuildDate, build_version.CommitHash),
	Args:         cobra.NoArgs,
	SilenceUsage: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(auth.GetCommand())
	rootCmd.AddCommand(send.GetCommand())
	rootCmd.AddCommand(serve.GetCommand())
}
