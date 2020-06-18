// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"io"

	"github.com/spf13/pflag"
)

const (
	CertificateFileFlag = "cert-file"
	CertificateFileDesc = "path to .p12 file containing APNs certificate"

	CertificatePasswordFlag = "cert-password"
	CertificatePasswordDesc = "password to decrypt the .p12 file (optional)"
)

type CertificateAuth struct {
	CertificateFile     string
	CertificatePassword string
	CertificateReader   io.Reader
}

func BindCertificateAuthFlags(flags *pflag.FlagSet, certificateAuth *CertificateAuth) {
	flags.StringVar(&certificateAuth.CertificateFile, CertificateFileFlag, certificateAuth.CertificateFile, CertificateFileDesc)
	flags.StringVar(&certificateAuth.CertificatePassword, CertificatePasswordFlag, certificateAuth.CertificatePassword, CertificatePasswordDesc)
}
