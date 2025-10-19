package publishing

import (
	"devopstoolkit/youtube-automation/internal/constants"
	"devopstoolkit/youtube-automation/internal/storage"
	"google.golang.org/api/youtube/v3"
)

// ValidateAndSetLanguage validates the language and sets it in the YouTube video object.
// It implements proper error handling with fallback mechanisms.
func ValidateAndSetLanguage(youtubeVideo *youtube.Video, video *storage.Video, defaultLanguage string) error {
	// Get the language to use (from video metadata or fallback to default)
	language := video.GetLanguage(defaultLanguage)
	audioLanguage := video.GetAudioLanguage(defaultLanguage)

	// Increment validation counter
	YouTubeMetrics.IncLanguageValidation()

	// Validate language codes
	if !constants.IsValidLanguage(language) {
		LogYouTubeWarn("Invalid language code '%s', falling back to default '%s'", language, defaultLanguage)
		YouTubeMetrics.IncLanguageFallback()
		language = defaultLanguage
	}

	if !constants.IsValidLanguage(audioLanguage) {
		LogYouTubeWarn("Invalid audio language code '%s', falling back to default '%s'", audioLanguage, defaultLanguage)
		YouTubeMetrics.IncLanguageFallback()
		audioLanguage = defaultLanguage
	}

	// Set language in video object with error handling
	err := setLanguageSafely(youtubeVideo, language, audioLanguage)
	if err != nil {
		// Log the error but don't fail the upload
		LogLanguageSetting(language, false, true, err)
		YouTubeMetrics.IncLanguageSetFailure()
		
		// Fallback to default language
		fallbackErr := setLanguageSafely(youtubeVideo, defaultLanguage, defaultLanguage)
		if fallbackErr != nil {
			// If even fallback fails, log but continue
			LogYouTubeError(NewLanguageError(defaultLanguage, fallbackErr), "Failed to set fallback language")
			YouTubeMetrics.IncLanguageSetFailure()
		} else {
			LogLanguageSetting(defaultLanguage, true, true, nil)
			YouTubeMetrics.IncLanguageSetSuccess()
		}
	} else {
		LogLanguageSetting(language, true, false, nil)
		YouTubeMetrics.IncLanguageSetSuccess()
	}

	// Store the applied languages back to the video struct
	video.AppliedLanguage = language
	video.AppliedAudioLanguage = audioLanguage

	return nil // Never fail the upload due to language setting issues
}

// setLanguageSafely sets the language fields on the YouTube video object.
// It handles potential nil pointer issues and other edge cases.
func setLanguageSafely(youtubeVideo *youtube.Video, language, audioLanguage string) error {
	if youtubeVideo == nil {
		return NewLanguageError(language, nil)
	}

	if youtubeVideo.Snippet == nil {
		// Create snippet if it doesn't exist
		youtubeVideo.Snippet = &youtube.VideoSnippet{}
	}

	// Set the language fields
	youtubeVideo.Snippet.DefaultLanguage = language
	youtubeVideo.Snippet.DefaultAudioLanguage = audioLanguage

	return nil
}

// ValidateLanguageCode validates a single language code and returns an error if invalid.
func ValidateLanguageCode(language string) error {
	if !constants.IsValidLanguage(language) {
		return NewLanguageError(language, nil)
	}
	return nil
}

// GetLanguageWithFallback returns the language to use with proper fallback logic.
func GetLanguageWithFallback(video *storage.Video, defaultLanguage string) (string, string) {
	language := video.GetLanguage(defaultLanguage)
	audioLanguage := video.GetAudioLanguage(defaultLanguage)

	// Validate and fallback if necessary
	if !constants.IsValidLanguage(language) {
		LogYouTubeWarn("Invalid language code '%s', using fallback '%s'", language, defaultLanguage)
		YouTubeMetrics.IncLanguageFallback()
		language = defaultLanguage
	}

	if !constants.IsValidLanguage(audioLanguage) {
		LogYouTubeWarn("Invalid audio language code '%s', using fallback '%s'", audioLanguage, defaultLanguage)
		YouTubeMetrics.IncLanguageFallback()
		audioLanguage = defaultLanguage
	}

	return language, audioLanguage
}
