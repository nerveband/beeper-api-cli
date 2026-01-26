package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorCategory represents the type of error for categorization
type ErrorCategory string

const (
	// CategoryAuth indicates authentication/authorization errors
	CategoryAuth ErrorCategory = "auth"
	// CategoryConfig indicates configuration errors
	CategoryConfig ErrorCategory = "config"
	// CategoryPermission indicates permission/scope errors
	CategoryPermission ErrorCategory = "permission"
	// CategoryNotFound indicates resource not found errors
	CategoryNotFound ErrorCategory = "not_found"
	// CategoryNetwork indicates network/connectivity errors
	CategoryNetwork ErrorCategory = "network"
	// CategoryValidation indicates input validation errors
	CategoryValidation ErrorCategory = "validation"
	// CategoryServer indicates server-side errors
	CategoryServer ErrorCategory = "server"
	// CategoryUnknown indicates unknown/uncategorized errors
	CategoryUnknown ErrorCategory = "unknown"
)

// APIError represents a structured error from the API
type APIError struct {
	// Message is the human-readable error message
	Message string `json:"message"`
	// Code is the error code (e.g., "auth_required", "chat_not_found")
	Code string `json:"code,omitempty"`
	// Category is the error category for hint generation
	Category ErrorCategory `json:"category"`
	// StatusCode is the HTTP status code
	StatusCode int `json:"status_code,omitempty"`
	// Operation is the operation that failed (e.g., "list_chats", "send_message")
	Operation string `json:"operation,omitempty"`
	// Hint is an actionable suggestion for resolving the error
	Hint string `json:"hint,omitempty"`
	// Underlying is the underlying error (not serialized to JSON)
	Underlying error `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *APIError) Unwrap() error {
	return e.Underlying
}

// ToJSON returns the error as a JSON string
func (e *APIError) ToJSON() string {
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{"message":%q,"category":"unknown"}`, e.Message)
	}
	return string(data)
}

// WithOperation sets the operation field and returns the error
func (e *APIError) WithOperation(op string) *APIError {
	e.Operation = op
	return e
}

// WithHint sets the hint field and returns the error
func (e *APIError) WithHint(hint string) *APIError {
	e.Hint = hint
	return e
}

// NewAPIError creates a new APIError with the given parameters
func NewAPIError(message string, category ErrorCategory) *APIError {
	return &APIError{
		Message:  message,
		Category: category,
	}
}

// NewAPIErrorFromStatus creates an APIError from an HTTP status code and response body
func NewAPIErrorFromStatus(statusCode int, body []byte, operation string) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		Operation:  operation,
	}

	// Try to parse error response from body
	var errResp struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Code    string `json:"code"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil {
		if errResp.Error != "" {
			apiErr.Message = errResp.Error
		} else if errResp.Message != "" {
			apiErr.Message = errResp.Message
		}
		apiErr.Code = errResp.Code
	}

	// Set default message if not parsed
	if apiErr.Message == "" {
		apiErr.Message = string(body)
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(statusCode)
		}
	}

	// Categorize based on status code
	apiErr.Category = categorizeStatusCode(statusCode)

	// Generate hint based on category and context
	apiErr.Hint = generateHint(apiErr)

	return apiErr
}

// categorizeStatusCode returns the appropriate error category for an HTTP status code
func categorizeStatusCode(statusCode int) ErrorCategory {
	switch {
	case statusCode == http.StatusUnauthorized:
		return CategoryAuth
	case statusCode == http.StatusForbidden:
		return CategoryPermission
	case statusCode == http.StatusNotFound:
		return CategoryNotFound
	case statusCode == http.StatusBadRequest:
		return CategoryValidation
	case statusCode >= 500:
		return CategoryServer
	default:
		return CategoryUnknown
	}
}

// generateHint generates an actionable hint based on error context
func generateHint(apiErr *APIError) string {
	switch apiErr.Category {
	case CategoryAuth:
		return "Set BEEPER_TOKEN environment variable with a valid API token. Generate one in Beeper Desktop settings."
	case CategoryPermission:
		return "Your token may lack the required scope. Check token permissions in Beeper Desktop settings."
	case CategoryNotFound:
		if apiErr.Operation == "get_chat" || apiErr.Operation == "list_messages" {
			return "Verify the chat ID is correct. Use 'beeper chats list' to see available chats."
		}
		return "The requested resource was not found. Verify the ID is correct."
	case CategoryNetwork:
		return "Check that Beeper Desktop is running and the API is enabled. Try 'beeper discover' to find the API."
	case CategoryConfig:
		return "Check your configuration with 'beeper config show'. Reset with 'beeper config set-url http://localhost:39867'."
	case CategoryValidation:
		return "Check the command arguments. Use --help for usage information."
	case CategoryServer:
		return "The Beeper Desktop API returned a server error. Try restarting Beeper Desktop."
	default:
		return ""
	}
}

// WrapNetworkError wraps a network error with appropriate context
func WrapNetworkError(err error, operation string) *APIError {
	return &APIError{
		Message:    fmt.Sprintf("failed to connect to API: %v", err),
		Category:   CategoryNetwork,
		Operation:  operation,
		Hint:       "Check that Beeper Desktop is running and the API is enabled. Try 'beeper discover' to find the API.",
		Underlying: err,
	}
}

// WrapConfigError wraps a configuration error with appropriate context
func WrapConfigError(err error, message string) *APIError {
	return &APIError{
		Message:    message,
		Category:   CategoryConfig,
		Hint:       "Check your configuration with 'beeper config show'. Edit ~/.beeper-api-cli/config.yaml if needed.",
		Underlying: err,
	}
}

// IsAuthError returns true if the error is an authentication error
func IsAuthError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Category == CategoryAuth
	}
	return false
}

// IsNetworkError returns true if the error is a network error
func IsNetworkError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Category == CategoryNetwork
	}
	return false
}

// IsNotFoundError returns true if the error is a not found error
func IsNotFoundError(err error) bool {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.Category == CategoryNotFound
	}
	return false
}
