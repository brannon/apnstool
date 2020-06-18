// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/brannon/apnstool/apns"
	"github.com/brannon/apnstool/operation"
)

const MaxFormMemory = 64 * 1024 * 1024

// type SendApiModel struct {
// 	AppId       string `json:"appId"`
// 	DeviceToken string `json:"deviceToken"`
// }

// type SendAlertApiModel struct {
// 	SendApiModel

// 	AlertText string `json:"alertText"`
// }

func (cmd *ServeCmd) handlePostSendAlert(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseMultipartForm(MaxFormMemory)
	if err != nil {
		WriteHtmlTemplate(rw, ErrorHtmlTemplate, err)
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
		WriteHtmlTemplate(rw, ErrorHtmlTemplate, errors.New("Invalid auth-type"))
		return
	}

	op.AlertText = getFormString(req, "alert-text")
	op.BadgeCount = getFormInt(req, "badge-count")
	op.SoundName = getFormString(req, "sound-name")

	fmt.Printf("OP: %+v\n", op)

	result, err := op.Exec()
	if err != nil {
		WriteHtmlTemplate(rw, ErrorHtmlTemplate, err)
		return
	}

	WriteHtmlTemplate(rw, ResultHtmlTemplate, result)
}

func getFormBool(req *http.Request, name string) bool {
	value := req.FormValue(name)
	if value != "" {
		boolValue, err := strconv.ParseBool(value)
		if err == nil {
			return boolValue
		}
	}

	return false
}

func getFormDurationOrDefault(req *http.Request, name string, defaultValue time.Duration) time.Duration {
	value := req.FormValue(name)
	if value != "" {
		durationValue, err := time.ParseDuration(value)
		if err == nil {
			return durationValue
		}
	}

	return defaultValue
}

func getFormFileReader(req *http.Request, name string) io.ReadCloser {
	file, _, err := req.FormFile(name)
	if err != nil {
		return nil
	}

	return file
}

func getFormInt(req *http.Request, name string) int {
	value := req.FormValue(name)
	if value != "" {
		intValue, err := strconv.ParseInt(value, 10, 32)
		if err == nil {
			return int(intValue)
		}
	}
	return -1
}

func getFormString(req *http.Request, name string) string {
	return req.FormValue(name)
}
