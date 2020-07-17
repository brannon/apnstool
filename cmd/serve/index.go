// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"net/http"
)

func (cmd *ServeCmd) handleIndex(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		switch req.URL.Path {
		case "/":
			WriteHtmlView(rw, req, "index", nil)
			return
		default:
			WriteNotFound(rw)
			return
		}
	default:
		WriteMethodNotAllowed(rw, []string{http.MethodGet})
		return
	}
}
