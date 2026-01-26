package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nerveband/beeper-api-cli/internal/api"
)

// JSONError represents an error in JSON format for --json-errors output
type JSONError struct {
	Error     string `json:"error"`
	Code      string `json:"code,omitempty"`
	Category  string `json:"category,omitempty"`
	Operation string `json:"operation,omitempty"`
	Hint      string `json:"hint,omitempty"`
}

// formatError formats an error for display, adding hints if available
func formatError(err error) string {
	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		var sb strings.Builder
		sb.WriteString("Error: ")
		sb.WriteString(apiErr.Message)
		sb.WriteString("\n")

		if apiErr.Hint != "" && !quietMode {
			sb.WriteString("\nHint: ")
			sb.WriteString(apiErr.Hint)
			sb.WriteString("\n")
		}

		return sb.String()
	}

	// For non-API errors, try to provide helpful hints based on error content
	errMsg := err.Error()
	hint := getGenericHint(errMsg)

	if hint != "" && !quietMode {
		return fmt.Sprintf("Error: %s\n\nHint: %s\n", errMsg, hint)
	}

	return fmt.Sprintf("Error: %s\n", errMsg)
}

// formatErrorAsJSON formats an error as JSON for --json-errors output
func formatErrorAsJSON(err error) string {
	jsonErr := JSONError{}

	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		jsonErr.Error = apiErr.Message
		jsonErr.Code = apiErr.Code
		jsonErr.Category = string(apiErr.Category)
		jsonErr.Operation = apiErr.Operation
		jsonErr.Hint = apiErr.Hint
	} else {
		jsonErr.Error = err.Error()
		jsonErr.Category = "unknown"
		jsonErr.Hint = getGenericHint(err.Error())
	}

	data, _ := json.Marshal(jsonErr)
	return string(data)
}

// getGenericHint returns a hint based on common error patterns
func getGenericHint(errMsg string) string {
	errLower := strings.ToLower(errMsg)

	switch {
	case strings.Contains(errLower, "connection refused"):
		return "Beeper Desktop may not be running. Start Beeper Desktop and ensure the API is enabled."
	case strings.Contains(errLower, "no such host"):
		return "Could not resolve the API host. Check your network connection and API URL configuration."
	case strings.Contains(errLower, "timeout"):
		return "The request timed out. Check if Beeper Desktop is responding and your network is stable."
	case strings.Contains(errLower, "unauthorized") || strings.Contains(errLower, "401"):
		return "Authentication required. Set BEEPER_TOKEN environment variable with a valid API token."
	case strings.Contains(errLower, "forbidden") || strings.Contains(errLower, "403"):
		return "Access denied. Your token may lack the required permissions."
	case strings.Contains(errLower, "not found") || strings.Contains(errLower, "404"):
		return "Resource not found. Verify the ID is correct using 'beeper chats list'."
	case strings.Contains(errLower, "config"):
		return "Configuration error. Check your settings with 'beeper config show'."
	default:
		return ""
	}
}

// handleError processes an error and exits with appropriate code
// Returns the exit code (for testing purposes)
func handleError(err error) int {
	if err == nil {
		return 0
	}

	// Determine exit code based on error type
	exitCode := 1 // Default: user/application error

	var apiErr *api.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Category {
		case api.CategoryNetwork:
			exitCode = 2 // System/network error
		case api.CategoryServer:
			exitCode = 2 // Server error (system)
		}
	} else {
		// Check for network-related errors in message
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "connection refused") ||
			strings.Contains(errMsg, "no such host") ||
			strings.Contains(errMsg, "timeout") {
			exitCode = 2
		}
	}

	// Output error
	if jsonErrors {
		fmt.Fprintln(os.Stderr, formatErrorAsJSON(err))
	} else {
		fmt.Fprint(os.Stderr, formatError(err))
	}

	return exitCode
}

// exitWithError handles an error and exits the process
func exitWithError(err error) {
	code := handleError(err)
	os.Exit(code)
}

// ExitCodes documents the exit codes used by the CLI
var ExitCodes = map[int]string{
	0: "Success",
	1: "User/Application error (invalid arguments, missing resources, permission denied)",
	2: "System/Network error (connection failed, timeout, server error)",
}
