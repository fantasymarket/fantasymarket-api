package responses

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error error `json:"error"`
}

type error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// CustomResponse encodes a custom message and sends a HTTP json response
func CustomResponse(w http.ResponseWriter, response interface{}, status int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// ErrorResponse writes a HTTP error with a json error message
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse{Error: error{
		Status:  statusCode,
		Message: message,
	}})
}
