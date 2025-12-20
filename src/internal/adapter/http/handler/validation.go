package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// validateRequest validates a DTO struct and returns an error if validation fails
func validateRequest(req interface{}) error {
	if err := validate.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			tag := err.Tag()
			var message string
			switch tag {
			case "required":
				message = fmt.Sprintf("%s is required", field)
			case "email":
				message = fmt.Sprintf("%s must be a valid email address", field)
			case "min":
				message = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
			case "max":
				message = fmt.Sprintf("%s must be at most %s characters", field, err.Param())
			case "oneof":
				message = fmt.Sprintf("%s must be one of: %s", field, err.Param())
			case "uuid":
				message = fmt.Sprintf("%s must be a valid UUID", field)
			case "url":
				message = fmt.Sprintf("%s must be a valid URL", field)
			case "gte":
				message = fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
			case "omitempty":
				// This should not trigger an error, but handle it just in case
				continue
			default:
				message = fmt.Sprintf("%s failed validation: %s", field, tag)
			}
			validationErrors = append(validationErrors, message)
		}
		return fmt.Errorf("validation failed: %s", strings.Join(validationErrors, "; "))
	}
	return nil
}

// validateAndRespond validates a request and responds with 400 if validation fails
func validateAndRespond(w http.ResponseWriter, req interface{}) bool {
	if err := validateRequest(req); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

