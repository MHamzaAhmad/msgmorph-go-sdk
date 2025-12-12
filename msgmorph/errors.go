package msgmorph

import (
	"encoding/json"
	"fmt"
)

// ErrorCode represents error codes returned by the MsgMorph API.
type ErrorCode string

// Error codes for MsgMorph API errors.
const (
	// Client errors
	ErrInvalidAPIKey        ErrorCode = "INVALID_API_KEY"
	ErrInvalidOrganizationID ErrorCode = "INVALID_ORGANIZATION_ID"
	ErrMissingRequiredField ErrorCode = "MISSING_REQUIRED_FIELD"
	ErrValidationError      ErrorCode = "VALIDATION_ERROR"

	// Authentication errors
	ErrUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrForbidden    ErrorCode = "FORBIDDEN"

	// Resource errors
	ErrNotFound      ErrorCode = "NOT_FOUND"
	ErrConflict      ErrorCode = "CONFLICT"
	ErrAlreadyExists ErrorCode = "ALREADY_EXISTS"

	// Server errors
	ErrInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// Network errors
	ErrNetworkError ErrorCode = "NETWORK_ERROR"
	ErrTimeout      ErrorCode = "TIMEOUT"
)

// errorMessages provides human-readable hints for common error codes.
var errorMessages = map[ErrorCode]string{
	ErrInvalidAPIKey:        "Invalid API key. Please check your MSGMORPH_API_KEY environment variable.",
	ErrInvalidOrganizationID: "Invalid organization ID. Please check your MSGMORPH_ORGANIZATION_ID environment variable.",
	ErrUnauthorized:         "Authentication failed. Please verify your API key is correct and has not expired.",
	ErrForbidden:            "Access denied. Your API key does not have permission to perform this action.",
	ErrNotFound:             "The requested resource was not found.",
	ErrConflict:             "A conflict occurred. The resource may already exist or be in an invalid state.",
	ErrAlreadyExists:        "This resource already exists. Use update instead of create.",
	ErrValidationError:      "Invalid request data. Please check the required fields.",
	ErrInternalError:        "An internal server error occurred. Please try again later.",
	ErrServiceUnavailable:   "The MsgMorph API is temporarily unavailable. Please try again later.",
	ErrNetworkError:         "Network error. Please check your internet connection and that the API URL is correct.",
	ErrTimeout:              "Request timed out. Please try again.",
}

// Error represents an error returned by the MsgMorph API.
//
// Example usage:
//
//	contact, err := client.Contacts.Create(ctx, input)
//	if err != nil {
//	    var msgErr *msgmorph.Error
//	    if errors.As(err, &msgErr) {
//	        fmt.Printf("Error: %s\n", msgErr.Message)
//	        fmt.Printf("Code: %s\n", msgErr.Code)
//	        fmt.Printf("Status: %d\n", msgErr.Status)
//	        fmt.Printf("Hint: %s\n", msgErr.Hint)
//	    }
//	}
type Error struct {
	// Message is the error message returned by the API.
	Message string `json:"message"`

	// Status is the HTTP status code.
	Status int `json:"status"`

	// Code is the error code that identifies the type of error.
	Code ErrorCode `json:"code"`

	// Hint provides a human-readable suggestion for resolving the error.
	Hint string `json:"hint,omitempty"`

	// Details contains additional error information.
	Details map[string]interface{} `json:"details,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Hint != "" && e.Hint != e.Message {
		return fmt.Sprintf("MsgMorphError [%s]: %s (Hint: %s)", e.Code, e.Message, e.Hint)
	}
	return fmt.Sprintf("MsgMorphError [%s]: %s", e.Code, e.Message)
}

// newError creates a new Error with the given parameters.
func newError(message string, status int, code ErrorCode, details map[string]interface{}) *Error {
	hint := errorMessages[code]
	if message == "" || message == "An unexpected error occurred" {
		if hint != "" {
			message = hint
		}
	}

	return &Error{
		Message: message,
		Status:  status,
		Code:    code,
		Hint:    hint,
		Details: details,
	}
}

// errorCodeFromStatus maps HTTP status codes to error codes.
func errorCodeFromStatus(status int) ErrorCode {
	switch status {
	case 400:
		return ErrValidationError
	case 401:
		return ErrUnauthorized
	case 403:
		return ErrForbidden
	case 404:
		return ErrNotFound
	case 409:
		return ErrConflict
	case 503:
		return ErrServiceUnavailable
	default:
		if status >= 500 {
			return ErrInternalError
		}
		return ErrValidationError
	}
}

// newNetworkError creates an Error for network-related failures.
func newNetworkError(err error) *Error {
	message := "Network request failed"
	if err != nil {
		message = err.Error()
	}

	return newError(message, 0, ErrNetworkError, nil)
}

// ToJSON converts the error to a JSON string for logging.
func (e *Error) ToJSON() string {
	data, _ := json.Marshal(e)
	return string(data)
}

// IsNotFound returns true if the error is a not found error.
func (e *Error) IsNotFound() bool {
	return e.Code == ErrNotFound
}

// IsUnauthorized returns true if the error is an authentication error.
func (e *Error) IsUnauthorized() bool {
	return e.Code == ErrUnauthorized
}

// IsValidationError returns true if the error is a validation error.
func (e *Error) IsValidationError() bool {
	return e.Code == ErrValidationError
}

// IsServerError returns true if the error is a server-side error.
func (e *Error) IsServerError() bool {
	return e.Code == ErrInternalError || e.Code == ErrServiceUnavailable
}
