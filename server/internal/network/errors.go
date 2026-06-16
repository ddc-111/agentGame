package network

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func NotFound(resource string) APIError {
	return APIError{Code: "NOT_FOUND", Message: resource + " not found"}
}

func BadRequest(msg string) APIError {
	return APIError{Code: "BAD_REQUEST", Message: msg}
}

func InternalError() APIError {
	return APIError{Code: "INTERNAL", Message: "internal server error"}
}

func respondError(c *gin.Context, status int, err APIError) {
	c.JSON(status, gin.H{"error": err})
}

func respondInternalError(c *gin.Context, actualErr error) {
	LoggerFromContext(c.Request.Context()).Error("Internal error", "error", actualErr)
	c.JSON(http.StatusInternalServerError, gin.H{"error": InternalError()})
}

func jsonMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func jsonUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
