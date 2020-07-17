// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package operation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/cmd/auth"
)

type SendOperationResult struct {
	ApnsId string
}

type SendOperationError struct {
	StatusCode  int
	ErrorReason string
}

func (err *SendOperationError) Error() string {
	return fmt.Sprintf("Operation failed with error: %d %s", err.StatusCode, err.ErrorReason)
}

func GetErrorStatusCode(err error) int {
	if sendOperationError, ok := err.(*SendOperationError); ok {
		return sendOperationError.StatusCode
	}
	return -1
}

type SendOperation struct {
	AppId           string
	CertificateAuth auth.CertificateAuth
	DeviceToken     string
	Sandbox         bool
	TokenAuth       auth.TokenAuth

	Client apns.Client
}

func (op *SendOperation) sendNotification(
	headers apns.Headers,
	content []byte,
) (*SendOperationResult, error) {
	if op.Sandbox {
		op.Client.ConfigureEndpoint(apns.SandboxEndpoint)
	}

	if op.useTokenAuth() {
		if op.TokenAuth.KeyReader == nil {
			keyFile, err := os.Open(op.TokenAuth.KeyFile)
			if err != nil {
				return nil, err
			}
			defer keyFile.Close()

			op.TokenAuth.KeyReader = keyFile
		}

		token, err := apns.GenerateJWTFromKeyReader(
			op.TokenAuth.KeyReader,
			op.TokenAuth.KeyId,
			op.TokenAuth.TeamId,
			time.Now(),
			op.TokenAuth.ExpiresAfter,
		)
		if err != nil {
			return nil, err
		}

		op.Client.ConfigureTokenAuth(token)
	} else if op.useCertificateAuth() {
		if op.CertificateAuth.CertificateReader == nil {
			certFile, err := os.Open(op.CertificateAuth.CertificateFile)
			if err != nil {
				return nil, err
			}
			defer certFile.Close()

			op.CertificateAuth.CertificateReader = certFile
		}

		cert, err := apns.LoadCertificateFromReader(op.CertificateAuth.CertificateReader, op.CertificateAuth.CertificatePassword)
		if err != nil {
			return nil, err
		}

		op.Client.ConfigureCertificateAuth(cert)
	}

	result, err := op.Client.Send(op.DeviceToken, headers, content)
	if err != nil {
		return nil, err
	}

	if result.Success() {
		return &SendOperationResult{
			ApnsId: result.Id(),
		}, nil
	}

	return nil, &SendOperationError{
		StatusCode:  result.StatusCode,
		ErrorReason: result.ErrorReason(),
	}
}

func (op *SendOperation) useCertificateAuth() bool {
	return (op.CertificateAuth.CertificateFile != "" || op.CertificateAuth.CertificateReader != nil)
}

func (op *SendOperation) useTokenAuth() bool {
	return (op.TokenAuth.KeyFile != "" || op.TokenAuth.KeyReader != nil) &&
		op.TokenAuth.KeyId != "" &&
		op.TokenAuth.TeamId != ""
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
