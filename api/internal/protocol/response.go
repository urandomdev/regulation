package protocol

import "fmt"

// ErrorResponse represents an error response
// with a code, message, and optional metadata.
//
//nolint:errname // ErrorResponse is a convention for error responses
type ErrorResponse struct {
	Code    string         `json:"code" cbor:"code"`
	Message string         `json:"message" cbor:"message"`
	Meta    map[string]any `json:"meta,omitempty" cbor:"meta,omitempty"`
}

func (r ErrorResponse) String() string {
	return fmt.Sprintf("%s: %s", r.Code, r.Message)
}

func (r ErrorResponse) Error() string {
	return r.String()
}

// Unauthorized returns an unauthorized error
func Unauthorized(message string) ErrorResponse {
	return ErrorResponse{
		Code:    UnauthorizedError,
		Message: message,
	}
}

// InvalidCredentials returns an invalid credentials error
func InvalidCredentials() ErrorResponse {
	return ErrorResponse{
		Code:    InvalidCredentialsError,
		Message: "invalid email or password",
	}
}
