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
)

type LayoutContext struct {
	ViewName    string
	ViewContext interface{}
}

func WriteError(rw http.ResponseWriter, statusCode int, err error) {
	WritePlainText(rw, http.StatusInternalServerError, err.Error())
}

func WriteHtmlView(rw http.ResponseWriter, viewName string, context interface{}) {
	WriteHtmlLayout(rw, "main", &LayoutContext{
		ViewName:    viewName,
		ViewContext: context,
	})
}

func WriteHtmlLayout(rw http.ResponseWriter, layoutName string, context *LayoutContext) {
	tmpl, err := LookupLayoutTemplate(layoutName)
	if err != nil {
		WriteError(rw, http.StatusInternalServerError, err)
		return
	}

	WriteHtmlTemplate(rw, tmpl, context)
}

func WriteHtmlTemplate(rw http.ResponseWriter, tmpl *template.Template, context interface{}) {
	rw.Header().Set("Content-Type", "text/html")
	rw.WriteHeader(200)

	err := tmpl.Execute(rw, context)
	if err != nil {
		fmt.Fprintf(rw, "<!-- Error rendering template: %s -->", err)
	}
}

func WriteMethodNotAllowed(rw http.ResponseWriter, allowedMethods []string) {
	text := fmt.Sprintf("Methods allowed: %s\n", strings.Join(allowedMethods, ","))
	WritePlainText(rw, http.StatusMethodNotAllowed, text)
}

func WritePlainText(rw http.ResponseWriter, statusCode int, text string) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(statusCode)

	_, _ = io.WriteString(rw, text)
}
