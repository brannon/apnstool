package operation

import (
	"fmt"
	"time"

	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/cmd/auth"
)

type SendOperationResult struct {
	ApnsId string
}

type SendOperation struct {
	AppId           string
	CertificateAuth auth.CertificateAuth
	DeviceToken     string
	Sandbox         bool
	TokenAuth       auth.TokenAuth
	Verbose         bool

	Client apns.Client
}

func (op *SendOperation) sendNotification(
	headers apns.Headers,
	content []byte,
) (*SendOperationResult, error) {
	// if op.Verbose {
	// 	op.Client.EnableLogging(op.IO.Stdout())
	// }

	if op.Sandbox {
		op.Client.ConfigureEndpoint(apns.SandboxEndpoint)
	}

	if op.useTokenAuth() {
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

	return nil, fmt.Errorf("Send operation failed with error: %d %s", result.StatusCode, result.ErrorReason())
}

func (op *SendOperation) useCertificateAuth() bool {
	return (op.CertificateAuth.CertificateFile != "" || op.CertificateAuth.CertificateReader != nil)
}

func (op *SendOperation) useTokenAuth() bool {
	return (op.TokenAuth.KeyFile != "" || op.TokenAuth.KeyReader != nil) &&
		op.TokenAuth.KeyId != "" &&
		op.TokenAuth.TeamId != ""
}
