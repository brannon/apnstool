package serve

import (
	"encoding/json"
	"net/http"
)

type JSONObject map[string]interface{}
type JSONArray []interface{}

type ErrorApiModel struct {
	Message string `json:"message"`
}

func ReadJSON(req *http.Request, body interface{}) error {
	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	return decoder.Decode(&body)
}

func WriteError(rw http.ResponseWriter, statusCode int, err error) {
	model := &ErrorApiModel{
		Message: err.Error(),
	}

	WriteJSON(rw, 500, model)
}

func WriteJSON(rw http.ResponseWriter, statusCode int, body interface{}) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	encoder := json.NewEncoder(rw)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(body)
}

func WriteOperationResult(rw http.ResponseWriter, statusCode int, result string) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(statusCode)
	rw.Write([]byte(result))
}
