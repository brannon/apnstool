// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package serve

import (
	"net/http"
)

func (cmd *ServeCmd) handleSendRaw(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		cmd.handleSendRawGet(rw, req)
		return
	case http.MethodPost:
		cmd.handleSendRawPost(rw, req)
		return
	default:
		WriteMethodNotAllowed(rw, []string{http.MethodGet, http.MethodPost})
		return
	}
}

func (cmd *ServeCmd) handleSendRawGet(rw http.ResponseWriter, req *http.Request) {

}

func (cmd *ServeCmd) handleSendRawPost(rw http.ResponseWriter, req *http.Request) {

}
