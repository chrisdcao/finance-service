package response

import (
	"encoding/json"
	"net/http"
)

// Response struct to encapsulate the response format
type Response struct {
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// WriteJSONResponse writes the response in JSON format
func WriteJSONResponse(w http.ResponseWriter, statusCode int, err string, data interface{}, message string) {
	response := Response{
		Error:   err,
		Data:    data,
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(response)
}
