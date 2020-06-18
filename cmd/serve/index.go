package serve

import (
	"net/http"
)

func (cmd *ServeCmd) handleGetIndex(rw http.ResponseWriter, req *http.Request) {
	WriteHtmlTemplate(rw, IndexHtmlTemplate, nil)
}
