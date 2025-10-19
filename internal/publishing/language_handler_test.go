package publishing

import (
	"devopstoolkit/youtube-automation/internal/storage"
	"google.golang.org/api/youtube/v3"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateAndSetLanguage(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	tests := []struct {
		name              string
		video             *storage.Video
		defaultLanguage   string
		expectedLanguage  string
		expectedAudioLang string
		expectError       bool
	}{
		{
			name: "Valid language codes",
			video: &storage.Video{
				Language:      "en",
				AudioLanguage: "en",
			},
			defaultLanguage:   "en",
			expectedLanguage:  "en",
			expectedAudioLang: "en",
			expectError:       false,
		},
		{
			name: "Empty language codes with fallback",
			video: &storage.Video{
				Language:      "",
				AudioLanguage: "",
			},
			defaultLanguage:   "fr",
			expectedLanguage:  "fr",
			expectedAudioLang: "fr",
			expectError:       false,
		},
		{
			name: "Invalid language codes with fallback",
			video: &storage.Video{
				Language:      "invalid",
				AudioLanguage: "invalid",
			},
			defaultLanguage:   "es",
			expectedLanguage:  "es",
			expectedAudioLang: "es",
			expectError:       false,
		},
		{
			name: "Mixed valid and invalid",
			video: &storage.Video{
				Language:      "en",
				AudioLanguage: "invalid",
			},
			defaultLanguage:   "fr",
			expectedLanguage:  "en",
			expectedAudioLang: "fr",
			expectError:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset metrics for each test
			YouTubeMetrics.Reset()

			youtubeVideo := &youtube.Video{
				Snippet: &youtube.VideoSnippet{},
			}

			err := ValidateAndSetLanguage(youtubeVideo, tt.video, tt.defaultLanguage)

			// Should never fail due to language setting
			assert.NoError(t, err)

			// Check that language was set correctly
			assert.Equal(t, tt.expectedLanguage, youtubeVideo.Snippet.DefaultLanguage)
			assert.Equal(t, tt.expectedAudioLang, youtubeVideo.Snippet.DefaultAudioLanguage)

			// Check that applied languages were stored
			assert.Equal(t, tt.expectedLanguage, tt.video.AppliedLanguage)
			assert.Equal(t, tt.expectedAudioLang, tt.video.AppliedAudioLanguage)
		})
	}
}

func TestValidateAndSetLanguage_NilVideo(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	video := &storage.Video{
		Language:      "en",
		AudioLanguage: "en",
	}

	// Test with nil YouTube video
	err := ValidateAndSetLanguage(nil, video, "en")

	// Should not fail the upload
	assert.NoError(t, err)
}

func TestValidateAndSetLanguage_NilSnippet(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	youtubeVideo := &youtube.Video{} // No snippet
	video := &storage.Video{
		Language:      "en",
		AudioLanguage: "en",
	}

	err := ValidateAndSetLanguage(youtubeVideo, video, "en")

	// Should not fail the upload
	assert.NoError(t, err)

	// Check that snippet was created
	assert.NotNil(t, youtubeVideo.Snippet)
	assert.Equal(t, "en", youtubeVideo.Snippet.DefaultLanguage)
	assert.Equal(t, "en", youtubeVideo.Snippet.DefaultAudioLanguage)
}

func TestValidateLanguageCode(t *testing.T) {
	tests := []struct {
		name        string
		language    string
		expectError bool
	}{
		{"Valid language", "en", false},
		{"Valid language", "es", false},
		{"Valid language", "fr", false},
		{"Invalid language", "invalid", true},
		{"Empty language", "", true},
		{"Invalid format", "english", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLanguageCode(tt.language)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetLanguageWithFallback(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	tests := []struct {
		name              string
		video             *storage.Video
		defaultLanguage   string
		expectedLanguage  string
		expectedAudioLang string
	}{
		{
			name: "Valid languages",
			video: &storage.Video{
				Language:      "en",
				AudioLanguage: "en",
			},
			defaultLanguage:   "fr",
			expectedLanguage:  "en",
			expectedAudioLang: "en",
		},
		{
			name: "Empty languages with fallback",
			video: &storage.Video{
				Language:      "",
				AudioLanguage: "",
			},
			defaultLanguage:   "es",
			expectedLanguage:  "es",
			expectedAudioLang: "es",
		},
		{
			name: "Invalid languages with fallback",
			video: &storage.Video{
				Language:      "invalid",
				AudioLanguage: "invalid",
			},
			defaultLanguage:   "de",
			expectedLanguage:  "de",
			expectedAudioLang: "de",
		},
		{
			name: "Mixed valid and invalid",
			video: &storage.Video{
				Language:      "en",
				AudioLanguage: "invalid",
			},
			defaultLanguage:   "fr",
			expectedLanguage:  "en",
			expectedAudioLang: "fr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset metrics for each test
			YouTubeMetrics.Reset()

			language, audioLanguage := GetLanguageWithFallback(tt.video, tt.defaultLanguage)

			assert.Equal(t, tt.expectedLanguage, language)
			assert.Equal(t, tt.expectedAudioLang, audioLanguage)
		})
	}
}

func TestValidateAndSetLanguage_Metrics(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	// Test with valid language
	video := &storage.Video{
		Language:      "en",
		AudioLanguage: "en",
	}
	youtubeVideo := &youtube.Video{
		Snippet: &youtube.VideoSnippet{},
	}

	err := ValidateAndSetLanguage(youtubeVideo, video, "en")
	assert.NoError(t, err)

	// Check metrics
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageValidation())
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageSetSuccess())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageSetFailure())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageFallback())
}

func TestValidateAndSetLanguage_InvalidLanguageMetrics(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	// Test with invalid language
	video := &storage.Video{
		Language:      "invalid",
		AudioLanguage: "invalid",
	}
	youtubeVideo := &youtube.Video{
		Snippet: &youtube.VideoSnippet{},
	}

	err := ValidateAndSetLanguage(youtubeVideo, video, "en")
	assert.NoError(t, err)

	// Check metrics
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageValidation())
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageSetSuccess())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageSetFailure())
	assert.Equal(t, int64(2), YouTubeMetrics.GetLanguageFallback()) // Both language and audio language fallback
}

func TestValidateAndSetLanguage_EdgeCases(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	tests := []struct {
		name        string
		video       *storage.Video
		expectError bool
	}{
		{
			name: "Nil video",
			video: nil,
			expectError: false,
		},
		{
			name: "Empty video",
			video: &storage.Video{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset metrics for each test
			YouTubeMetrics.Reset()

			youtubeVideo := &youtube.Video{
				Snippet: &youtube.VideoSnippet{},
			}

			err := ValidateAndSetLanguage(youtubeVideo, tt.video, "en")
			assert.NoError(t, err)
		})
	}
}
