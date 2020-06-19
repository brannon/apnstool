// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"io"
	"net/http"
	"strconv"
	"time"
)

const MaxFormMemory = 64 * 1024 * 1024

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
