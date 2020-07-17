package serve

import "net/http"

func (cmd *ServeCmd) handleFaq(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		WriteHtmlView(rw, req, "faq", nil)
		return

	default:
		WriteMethodNotAllowed(rw, []string{http.MethodGet})
		return
	}
}
