package network

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
