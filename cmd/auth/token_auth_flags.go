// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"time"

	"github.com/spf13/pflag"
)

const (
	KeyFileFlag = "key-file"
	KeyFileDesc = "path to .p8 file containing APNs-enabled private key"

	KeyIdFlag = "key-id"
	KeyIdDesc = "key ID for the APNs-enabled private key"

	TeamIdFlag = "team-id"
	TeamIdDesc = "Apple Developer team ID"

	ExpiresAfterFlag    = "expires-after"
	ExpiresAfterDefault = 60 * time.Minute
	ExpiresAfterDesc    = "amount of time until the token expires"
)

type TokenAuth struct {
	KeyFile      string
	KeyId        string
	TeamId       string
	ExpiresAfter time.Duration
}

func BindTokenAuthFlags(flags *pflag.FlagSet, tokenAuth *TokenAuth) {
	flags.StringVar(&tokenAuth.KeyFile, KeyFileFlag, tokenAuth.KeyFile, KeyFileDesc)
	flags.StringVar(&tokenAuth.KeyId, KeyIdFlag, tokenAuth.KeyId, KeyIdDesc)
	flags.StringVar(&tokenAuth.TeamId, TeamIdFlag, tokenAuth.TeamId, TeamIdDesc)
	flags.DurationVar(&tokenAuth.ExpiresAfter, ExpiresAfterFlag, ExpiresAfterDefault, ExpiresAfterDesc)
}
