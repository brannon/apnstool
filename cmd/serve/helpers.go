// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/brannon/apnstool/build_version"
)

type BuildContext struct {
	BuildDate  string
	CommitHash string
	Version    string
}

type RequestContext struct {
	Path string
}

type LayoutContext struct {
	Build   BuildContext
	Request RequestContext

	ViewName    string
	ViewContext interface{}
}

func WriteError(rw http.ResponseWriter, statusCode int, err error) {
	WritePlainText(rw, http.StatusInternalServerError, err.Error())
}

func WriteHtmlView(rw http.ResponseWriter, req *http.Request, viewName string, context interface{}) {
	WriteHtmlLayout(rw, http.StatusOK, "main", &LayoutContext{
		Build: BuildContext{
			BuildDate:  build_version.BuildDate,
			CommitHash: build_version.CommitHash,
		},
		Request: RequestContext{
			Path: req.URL.Path,
		},
		ViewName:    viewName,
		ViewContext: context,
	})
}

func WriteHtmlViewWithStatus(rw http.ResponseWriter, statusCode int, req *http.Request, viewName string, context interface{}) {
	WriteHtmlLayout(rw, statusCode, "main", &LayoutContext{
		Build: BuildContext{
			BuildDate:  build_version.BuildDate,
			CommitHash: build_version.CommitHash,
		},
		Request: RequestContext{
			Path: req.URL.Path,
		},
		ViewName:    viewName,
		ViewContext: context,
	})
}

func WriteHtmlLayout(rw http.ResponseWriter, statusCode int, layoutName string, context *LayoutContext) {
	tmpl, err := LookupLayoutTemplate(layoutName)
	if err != nil {
		WriteError(rw, http.StatusInternalServerError, err)
		return
	}

	WriteHtmlTemplate(rw, statusCode, tmpl, context)
}

func WriteHtmlTemplate(rw http.ResponseWriter, statusCode int, tmpl *template.Template, context interface{}) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(statusCode)

	err := tmpl.Execute(rw, context)
	if err != nil {
		fmt.Fprintf(rw, "<!-- Error rendering template: %s -->", err)
	}
}

func WriteMethodNotAllowed(rw http.ResponseWriter, allowedMethods []string) {
	text := fmt.Sprintf("Methods allowed: %s\n", strings.Join(allowedMethods, ","))
	WritePlainText(rw, http.StatusMethodNotAllowed, text)
}

func WriteNotFound(rw http.ResponseWriter) {
	WritePlainText(rw, http.StatusNotFound, "Not Found")
}

func WritePlainText(rw http.ResponseWriter, statusCode int, text string) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(statusCode)

	_, _ = io.WriteString(rw, text)
}
