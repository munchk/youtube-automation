package publishing

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategorizeError(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		expectedType   ErrorType
		expectedRetry  bool
		expectedMsg    string
	}{
		{
			name:           "Authentication error",
			inputError:     errors.New("authentication failed"),
			expectedType:   ErrorTypeAuth,
			expectedRetry:  false,
			expectedMsg:    "Authentication failed or insufficient permissions",
		},
		{
			name:           "Rate limit error",
			inputError:     errors.New("rate limit exceeded"),
			expectedType:   ErrorTypeRateLimit,
			expectedRetry:  true,
			expectedMsg:    "Rate limit exceeded or quota exceeded",
		},
		{
			name:           "Network error",
			inputError:     errors.New("network timeout"),
			expectedType:   ErrorTypeNetwork,
			expectedRetry:  true,
			expectedMsg:    "Network connectivity issue",
		},
		{
			name:           "Invalid request error",
			inputError:     errors.New("invalid request"),
			expectedType:   ErrorTypeInvalid,
			expectedRetry:  false,
			expectedMsg:    "Invalid request or malformed data",
		},
		{
			name:           "Server error",
			inputError:     errors.New("internal server error"),
			expectedType:   ErrorTypeServer,
			expectedRetry:  true,
			expectedMsg:    "YouTube server error",
		},
		{
			name:           "Language error",
			inputError:     errors.New("language setting failed"),
			expectedType:   ErrorTypeLanguage,
			expectedRetry:  false,
			expectedMsg:    "Language setting error",
		},
		{
			name:           "Upload error",
			inputError:     errors.New("video upload failed"),
			expectedType:   ErrorTypeUpload,
			expectedRetry:  true,
			expectedMsg:    "Video upload error",
		},
		{
			name:           "Unknown error",
			inputError:     errors.New("some random error"),
			expectedType:   ErrorTypeUnknown,
			expectedRetry:  false,
			expectedMsg:    "Unknown error occurred",
		},
		{
			name:           "Nil error",
			inputError:     nil,
			expectedType:   ErrorTypeUnknown,
			expectedRetry:  false,
			expectedMsg:    "Unknown error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CategorizeError(tt.inputError)

			assert.Equal(t, tt.expectedType, result.Type)
			assert.Equal(t, tt.expectedRetry, result.Retryable)
			assert.Equal(t, tt.expectedMsg, result.Message)
			assert.Equal(t, tt.inputError, result.OriginalError)
		})
	}
}

func TestNewLanguageError(t *testing.T) {
	originalErr := errors.New("original error")
	language := "en"

	langErr := NewLanguageError(language, originalErr)

	assert.Equal(t, ErrorTypeLanguage, langErr.Type)
	assert.Equal(t, language, langErr.Language)
	assert.Equal(t, originalErr, langErr.OriginalError)
	assert.False(t, langErr.Retryable)
	assert.Contains(t, langErr.Message, language)
}

func TestNewUploadError(t *testing.T) {
	originalErr := errors.New("upload failed")
	videoID := "test-video-123"

	uploadErr := NewUploadError(videoID, originalErr)

	assert.Equal(t, ErrorTypeUpload, uploadErr.Type)
	assert.Equal(t, videoID, uploadErr.VideoID)
	assert.Equal(t, originalErr, uploadErr.OriginalError)
	assert.True(t, uploadErr.Retryable)
	assert.Equal(t, "Video upload failed", uploadErr.Message)
}

func TestYouTubeError_Error(t *testing.T) {
	tests := []struct {
		name        string
		youtubeErr  *YouTubeError
		expectedMsg string
	}{
		{
			name: "Error with original error",
			youtubeErr: &YouTubeError{
				Type:          ErrorTypeAuth,
				Message:       "Authentication failed",
				OriginalError: errors.New("original error"),
			},
			expectedMsg: "YouTube error [auth]: Authentication failed (original: original error)",
		},
		{
			name: "Error without original error",
			youtubeErr: &YouTubeError{
				Type:    ErrorTypeLanguage,
				Message: "Language setting failed",
			},
			expectedMsg: "YouTube error [language_error]: Language setting failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.youtubeErr.Error()
			assert.Equal(t, tt.expectedMsg, result)
		})
	}
}

func TestYouTubeError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	youtubeErr := &YouTubeError{
		OriginalError: originalErr,
	}

	assert.Equal(t, originalErr, youtubeErr.Unwrap())
}

func TestCategorizeError_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
		expectedType ErrorType
	}{
		{"Uppercase auth", "AUTHENTICATION FAILED", ErrorTypeAuth},
		{"Mixed case rate limit", "Rate Limit Exceeded", ErrorTypeRateLimit},
		{"Lowercase network", "network timeout", ErrorTypeNetwork},
		{"Mixed case invalid", "Invalid Request", ErrorTypeInvalid},
		{"Uppercase server", "SERVER ERROR", ErrorTypeServer},
		{"Mixed case language", "Language Setting Failed", ErrorTypeLanguage},
		{"Uppercase upload", "VIDEO UPLOAD FAILED", ErrorTypeUpload},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errorMessage)
			result := CategorizeError(err)
			assert.Equal(t, tt.expectedType, result.Type)
		})
	}
}

func TestCategorizeError_MultipleKeywords(t *testing.T) {
	tests := []struct {
		name         string
		errorMessage string
		expectedType ErrorType
	}{
		{"Auth with unauthorized", "authentication failed: unauthorized", ErrorTypeAuth},
		{"Rate limit with quota", "rate limit exceeded: quota exceeded", ErrorTypeRateLimit},
		{"Network with connection", "network error: connection timeout", ErrorTypeNetwork},
		{"Invalid with bad request", "invalid request: bad request", ErrorTypeInvalid},
		{"Server with internal", "server error: internal server error", ErrorTypeServer},
		{"Language with locale", "language error: locale setting failed", ErrorTypeLanguage},
		{"Upload with video", "upload error: video processing failed", ErrorTypeUpload},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errors.New(tt.errorMessage)
			result := CategorizeError(err)
			assert.Equal(t, tt.expectedType, result.Type)
		})
	}
}
