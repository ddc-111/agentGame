package network

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	return v[0].Field + ": " + v[0].Message
}

func validateRequired(fields map[string]interface{}) ValidationErrors {
	var errs ValidationErrors
	for name, val := range fields {
		if val == nil {
			errs = append(errs, ValidationError{Field: name, Message: "required"})
			continue
		}
		if s, ok := val.(string); ok && s == "" {
			errs = append(errs, ValidationError{Field: name, Message: "required"})
		}
	}
	return errs
}

func validatePositiveInt(field string, val uint) ValidationErrors {
	var errs ValidationErrors
	if val == 0 {
		errs = append(errs, ValidationError{Field: field, Message: "must be a positive integer"})
	}
	return errs
}

func validateIntMin(field string, val, min int) ValidationErrors {
	var errs ValidationErrors
	if val < min {
		errs = append(errs, ValidationError{Field: field, Message: fmt.Sprintf("must be at least %d", min)})
	}
	return errs
}

func validateIntRange(field string, val, min, max int) ValidationErrors {
	var errs ValidationErrors
	if val < min || val > max {
		errs = append(errs, ValidationError{Field: field, Message: fmt.Sprintf("must be between %d and %d", min, max)})
	}
	return errs
}

func validateStringIn(field, val string, allowed []string) ValidationErrors {
	var errs ValidationErrors
	for _, a := range allowed {
		if val == a {
			return errs
		}
	}
	errs = append(errs, ValidationError{Field: field, Message: fmt.Sprintf("must be one of: %v", allowed)})
	return errs
}

func validateStringMaxLen(field, val string, maxLen int) ValidationErrors {
	var errs ValidationErrors
	if len(val) > maxLen {
		errs = append(errs, ValidationError{Field: field, Message: fmt.Sprintf("must be at most %d characters", maxLen)})
	}
	return errs
}

func mergeErrors(errs ...ValidationErrors) ValidationErrors {
	var result ValidationErrors
	for _, e := range errs {
		result = append(result, e...)
	}
	return result
}

func respondValidation(c *gin.Context, errs ValidationErrors) {
	if len(errs) == 0 {
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": APIError{
			Code:    "VALIDATION_ERROR",
			Message: "validation failed",
		},
		"validation_errors": errs,
	})
}
