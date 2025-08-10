package validation

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// APIError struct to encapsulate error details
type APIError struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Errors     map[string]string `json:"errors,omitempty"`
}

// Convert APIError to a string for error display
func (e *APIError) Error() string {
	err, _ := json.Marshal(e)
	return string(err)
}

// Factory function to create a new APIError
func NewAPIError(statusCode int, message string, validationErrors map[string]string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
		Errors:     validationErrors,
	}
}

// Reusable function to bind and validate a request struct
func BindAndValidate(c *fiber.Ctx, s interface{}) error {
	// Parse the request body into the given struct
	if err := c.BodyParser(s); err != nil {
		return NewAPIError(fiber.StatusBadRequest, err.Error(), nil)
	}

	// Validate the struct using the validator package
	validate := validator.New()

	// Register a custom tag name function to use the JSON tag name in errors.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" {
			return ""
		}
		// In case there are multiple options like "winner_count,omitempty"
		return strings.Split(name, ",")[0]
	})

	if err := validate.Struct(s); err != nil {
		// Collect validation errors
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", err.Field(), err.Tag())
		}
		// Return an APIError with validation errors
		return NewAPIError(fiber.StatusBadRequest, "Validation failed", validationErrors)
	}

	return nil
}
