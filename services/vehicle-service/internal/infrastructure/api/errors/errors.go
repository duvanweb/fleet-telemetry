package errors

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type errorResponse struct {
	Error string `json:"error"`
}

// WriteError writes a JSON error response with the given HTTP status code.
func WriteError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}
