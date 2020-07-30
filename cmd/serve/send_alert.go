// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"errors"
	"net/http"
	"time"

	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/operation"
)

func (cmd *ServeCmd) handleSendAlert(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		cmd.handleSendAlertGet(rw, req)
		return
	case http.MethodPost:
		cmd.handleSendAlertPost(rw, req)
		return
	default:
		WriteMethodNotAllowed(rw, []string{http.MethodGet, http.MethodPost})
		return
	}
}

func (cmd *ServeCmd) handleSendAlertGet(rw http.ResponseWriter, req *http.Request) {
	WriteHtmlView(rw, req, "send_alert", nil)
}

func (cmd *ServeCmd) handleSendAlertPost(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(MaxFormMemory)
	if err != nil {
		WriteHtmlView(rw, req, "error", err)
		return
	}

	appId := getFormString(req, "app-id")
	deviceToken := getFormString(req, "device-token")

	SetLoggingContextValue(rw, LogContextAppId, appId)
	SetLoggingContextValue(rw, LogContextDeviceToken, deviceToken)

	op := operation.NewSendAlert(apns.NewClient())

	op.AppId = appId
	op.DeviceToken = deviceToken
	op.Sandbox = getFormBool(req, "sandbox")

	authType := getFormString(req, "auth-type")
	if authType == "token" {
		op.TokenAuth.ExpiresAfter = getFormDurationOrDefault(req, "expires-after", 1*time.Hour)
		op.TokenAuth.KeyReader = getFormFileReader(req, "key-file")
		op.TokenAuth.KeyId = getFormString(req, "key-id")
		op.TokenAuth.TeamId = getFormString(req, "team-id")
	} else if authType == "cert" {
		op.CertificateAuth.CertificateReader = getFormFileReader(req, "cert-file")
		op.CertificateAuth.CertificatePassword = getFormString(req, "cert-password")
	} else {
		WriteHtmlView(rw, req, "error", errors.New("Invalid auth-type"))
		return
	}

	op.Expiration = getFormString(req, "expiration")
	op.Priority = getFormString(req, "priority")

	op.AlertText = getFormString(req, "alert-text")
	op.BadgeCount = getFormInt(req, "badge-count")
	op.SoundName = getFormString(req, "sound-name")

	result, err := op.Exec()
	if err != nil {
		statusCode := operation.GetErrorStatusCode(err)
		if statusCode == -1 {
			statusCode = http.StatusInternalServerError
		}
		WriteHtmlViewWithStatus(rw, statusCode, req, "error", err)
		return
	}

	WriteHtmlView(rw, req, "send_result", result)
}
