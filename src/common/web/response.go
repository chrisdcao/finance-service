package web

import (
	"github.com/gin-gonic/gin"
)

// Response struct to encapsulate the response format
type Response struct {
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// WriteJSONResponse writes the response in JSON format for Gin
func WriteJSONResponse(ctx *gin.Context, statusCode int, err string, data interface{}, message string) {
	response := Response{
		Error:   err,
		Data:    data,
		Message: message,
	}
	ctx.JSON(statusCode, response)
}
