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
	WriteHtmlView(rw, "send_alert", nil)
}

func (cmd *ServeCmd) handleSendAlertPost(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(MaxFormMemory)
	if err != nil {
		WriteHtmlView(rw, "error", err)
		return
	}

	op := operation.NewSendAlert(apns.NewClient())

	op.AppId = getFormString(req, "app-id")
	op.DeviceToken = getFormString(req, "device-token")
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
		WriteHtmlView(rw, "error", errors.New("Invalid auth-type"))
		return
	}

	op.AlertText = getFormString(req, "alert-text")
	op.BadgeCount = getFormInt(req, "badge-count")
	op.SoundName = getFormString(req, "sound-name")

	result, err := op.Exec()
	if err != nil {
		WriteHtmlView(rw, "error", err)
		return
	}

	WriteHtmlView(rw, "send_result", result)
}
