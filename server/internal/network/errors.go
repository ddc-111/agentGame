package network

import (
	"log"
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
	log.Printf("Internal error: %v", actualErr)
	c.JSON(http.StatusInternalServerError, gin.H{"error": InternalError()})
}
