package handler

import (
	"encoding/json"
	"net/http"
)

// Handler : function waiting for for an HTTP request, reads it and adecide what action to take
// Received values : http.ResponseWriter : used to write response back to client
// http.Request : contains all information about incoming request, such as URL, method, body, etc.
// Handlers connect web request to the right code

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}
