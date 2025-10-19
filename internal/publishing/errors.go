package publishing

import (
	"fmt"
	"strings"
)

// ErrorType defines the category of a YouTube-related error.
// This helps in deciding how to handle the error (e.g., retry, log level).
type ErrorType string

// YouTube API Error Categories
const (
	ErrorTypeAuth      ErrorType = "auth"            // Authentication or permission issue
	ErrorTypeRateLimit ErrorType = "rate_limit"      // Rate limit exceeded
	ErrorTypeNetwork   ErrorType = "network"         // Network connectivity problem
	ErrorTypeInvalid   ErrorType = "invalid_request"  // Malformed or invalid request
	ErrorTypeServer    ErrorType = "server_error"     // YouTube server-side issue (5xx errors)
	ErrorTypeLanguage  ErrorType = "language_error"  // Language setting specific errors
	ErrorTypeUpload    ErrorType = "upload_error"    // Video upload specific errors
	ErrorTypeUnknown   ErrorType = "unknown"         // Error that doesn't fit other categories
	ErrorTypeInternal  ErrorType = "internal"        // Errors originating from within this application
)

// YouTubeError is a custom error structure to wrap and categorize errors from YouTube operations.
type YouTubeError struct {
	Type          ErrorType // Category of the error
	Message       string    // Human-readable error message
	Retryable     bool      // Indicates if the operation that caused this error can be retried
	OriginalError error     // The original error object, if any
	VideoID       string    // Video ID if applicable
	Language      string    // Language code if applicable
}

// Error implements the error interface for YouTubeError.
func (e *YouTubeError) Error() string {
	if e.OriginalError != nil {
		return fmt.Sprintf("YouTube error [%s]: %s (original: %v)", e.Type, e.Message, e.OriginalError)
	}
	return fmt.Sprintf("YouTube error [%s]: %s", e.Type, e.Message)
}

// Unwrap provides compatibility for Go 1.13+ error chains.
func (e *YouTubeError) Unwrap() error {
	return e.OriginalError
}

// CategorizeError inspects an error and returns a structured YouTubeError.
// It attempts to identify specific error types from the YouTube API,
// then falls back to string matching for common error messages.
func CategorizeError(err error) *YouTubeError {
	if err == nil {
		return nil
	}

	// Fallback to string matching for common error patterns
	errStr := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errStr, "authentication") || strings.Contains(errStr, "unauthorized"):
		return &YouTubeError{
			Type:          ErrorTypeAuth,
			Message:       "Authentication failed or insufficient permissions",
			Retryable:     false,
			OriginalError: err,
		}
	case strings.Contains(errStr, "rate limit") || strings.Contains(errStr, "quota"):
		return &YouTubeError{
			Type:          ErrorTypeRateLimit,
			Message:       "Rate limit exceeded or quota exceeded",
			Retryable:     true,
			OriginalError: err,
		}
	case strings.Contains(errStr, "network") || strings.Contains(errStr, "timeout") || strings.Contains(errStr, "connection"):
		return &YouTubeError{
			Type:          ErrorTypeNetwork,
			Message:       "Network connectivity issue",
			Retryable:     true,
			OriginalError: err,
		}
	case strings.Contains(errStr, "invalid") || strings.Contains(errStr, "bad request"):
		return &YouTubeError{
			Type:          ErrorTypeInvalid,
			Message:       "Invalid request or malformed data",
			Retryable:     false,
			OriginalError: err,
		}
	case strings.Contains(errStr, "server error") || strings.Contains(errStr, "internal server"):
		return &YouTubeError{
			Type:          ErrorTypeServer,
			Message:       "YouTube server error",
			Retryable:     true,
			OriginalError: err,
		}
	case strings.Contains(errStr, "language") || strings.Contains(errStr, "locale"):
		return &YouTubeError{
			Type:          ErrorTypeLanguage,
			Message:       "Language setting error",
			Retryable:     false,
			OriginalError: err,
		}
	case strings.Contains(errStr, "upload") || strings.Contains(errStr, "video"):
		return &YouTubeError{
			Type:          ErrorTypeUpload,
			Message:       "Video upload error",
			Retryable:     true,
			OriginalError: err,
		}
	default:
		return &YouTubeError{
			Type:          ErrorTypeUnknown,
			Message:       "Unknown error occurred",
			Retryable:     false,
			OriginalError: err,
		}
	}
}

// NewLanguageError creates a specific error for language setting failures.
func NewLanguageError(language string, originalErr error) *YouTubeError {
	return &YouTubeError{
		Type:          ErrorTypeLanguage,
		Message:       fmt.Sprintf("Failed to set language to '%s'", language),
		Retryable:     false,
		OriginalError: originalErr,
		Language:      language,
	}
}

// NewUploadError creates a specific error for upload failures.
func NewUploadError(videoID string, originalErr error) *YouTubeError {
	return &YouTubeError{
		Type:          ErrorTypeUpload,
		Message:       "Video upload failed",
		Retryable:     true,
		OriginalError: originalErr,
		VideoID:       videoID,
	}
}
