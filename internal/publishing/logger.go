package publishing

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var youtubeLog *logrus.Logger

func init() {
	youtubeLog = logrus.New()
	youtubeLog.SetFormatter(&logrus.JSONFormatter{})
	// Default to Info level, can be made configurable later if needed
	youtubeLog.SetLevel(logrus.InfoLevel)
	youtubeLog.SetOutput(os.Stdout)
}

// SetLogLevel allows changing the global log level for YouTube operations.
func SetLogLevel(level logrus.Level) {
	youtubeLog.SetLevel(level)
}

func baseEntry() *logrus.Entry {
	return youtubeLog.WithField("component", "youtube")
}

// LogYouTubeError logs a categorized YouTube error with structured fields.
func LogYouTubeError(yErr *YouTubeError, message string) {
	if yErr == nil {
		baseEntry().Error(message)
		return
	}

	fields := logrus.Fields{
		"error_type": yErr.Type,
		"retryable":  yErr.Retryable,
	}
	
	// Add context fields if available
	if yErr.VideoID != "" {
		fields["video_id"] = yErr.VideoID
	}
	if yErr.Language != "" {
		fields["language"] = yErr.Language
	}

	entry := baseEntry().WithFields(fields)

	if yErr.OriginalError != nil {
		entry.WithError(yErr.OriginalError).Error(fmt.Sprintf("%s: %s", message, yErr.Message))
	} else {
		entry.Error(fmt.Sprintf("%s: %s", message, yErr.Message))
	}
}

// LogYouTubeWarn logs a warning message related to YouTube operations.
func LogYouTubeWarn(message string, args ...interface{}) {
	baseEntry().Warnf(message, args...)
}

// LogYouTubeInfo logs an informational message related to YouTube operations.
func LogYouTubeInfo(message string, args ...interface{}) {
	baseEntry().Infof(message, args...)
}

// LogYouTubeDebug logs a debug message related to YouTube operations.
func LogYouTubeDebug(message string, args ...interface{}) {
	baseEntry().Debugf(message, args...)
}

// LogLanguageSetting logs language setting operations with context.
func LogLanguageSetting(language string, success bool, fallback bool, err error) {
	fields := logrus.Fields{
		"language": language,
		"success":  success,
		"fallback": fallback,
	}

	entry := baseEntry().WithFields(fields)

	if err != nil {
		entry.WithError(err).Error("Language setting failed")
	} else if fallback {
		entry.Warn("Language setting succeeded with fallback to default")
	} else {
		entry.Info("Language setting succeeded")
	}
}

// LogUploadOperation logs upload operations with context.
func LogUploadOperation(videoID string, success bool, err error) {
	fields := logrus.Fields{
		"video_id": videoID,
		"success":  success,
	}

	entry := baseEntry().WithFields(fields)

	if err != nil {
		entry.WithError(err).Error("Upload operation failed")
	} else {
		entry.Info("Upload operation succeeded")
	}
}
