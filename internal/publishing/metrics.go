package publishing

import (
	"sync/atomic"
)

// Metrics tracks various YouTube operation statistics.
type Metrics struct {
	LanguageSetSuccess   int64 // Counter for successful language settings
	LanguageSetFailure   int64 // Counter for failed language settings
	UploadSuccess        int64 // Counter for successful uploads
	UploadFailure        int64 // Counter for failed uploads
	LanguageValidation   int64 // Counter for language validations
	LanguageFallback     int64 // Counter for language fallbacks to default
}

// YouTubeMetrics is the global metrics instance.
var YouTubeMetrics = &Metrics{}

// IncLanguageSetSuccess increments the successful language setting counter.
func (m *Metrics) IncLanguageSetSuccess() {
	atomic.AddInt64(&m.LanguageSetSuccess, 1)
}

// IncLanguageSetFailure increments the failed language setting counter.
func (m *Metrics) IncLanguageSetFailure() {
	atomic.AddInt64(&m.LanguageSetFailure, 1)
}

// IncUploadSuccess increments the successful upload counter.
func (m *Metrics) IncUploadSuccess() {
	atomic.AddInt64(&m.UploadSuccess, 1)
}

// IncUploadFailure increments the failed upload counter.
func (m *Metrics) IncUploadFailure() {
	atomic.AddInt64(&m.UploadFailure, 1)
}

// IncLanguageValidation increments the language validation counter.
func (m *Metrics) IncLanguageValidation() {
	atomic.AddInt64(&m.LanguageValidation, 1)
}

// IncLanguageFallback increments the language fallback counter.
func (m *Metrics) IncLanguageFallback() {
	atomic.AddInt64(&m.LanguageFallback, 1)
}

// GetLanguageSetSuccess returns the current value of successful language settings.
func (m *Metrics) GetLanguageSetSuccess() int64 {
	return atomic.LoadInt64(&m.LanguageSetSuccess)
}

// GetLanguageSetFailure returns the current value of failed language settings.
func (m *Metrics) GetLanguageSetFailure() int64 {
	return atomic.LoadInt64(&m.LanguageSetFailure)
}

// GetUploadSuccess returns the current value of successful uploads.
func (m *Metrics) GetUploadSuccess() int64 {
	return atomic.LoadInt64(&m.UploadSuccess)
}

// GetUploadFailure returns the current value of failed uploads.
func (m *Metrics) GetUploadFailure() int64 {
	return atomic.LoadInt64(&m.UploadFailure)
}

// GetLanguageValidation returns the current value of language validations.
func (m *Metrics) GetLanguageValidation() int64 {
	return atomic.LoadInt64(&m.LanguageValidation)
}

// GetLanguageFallback returns the current value of language fallbacks.
func (m *Metrics) GetLanguageFallback() int64 {
	return atomic.LoadInt64(&m.LanguageFallback)
}

// GetLanguageSetTotal returns the total number of language setting attempts.
func (m *Metrics) GetLanguageSetTotal() int64 {
	return m.GetLanguageSetSuccess() + m.GetLanguageSetFailure()
}

// GetUploadTotal returns the total number of upload attempts.
func (m *Metrics) GetUploadTotal() int64 {
	return m.GetUploadSuccess() + m.GetUploadFailure()
}

// GetLanguageSetSuccessRate returns the success rate for language setting (0.0 to 1.0).
func (m *Metrics) GetLanguageSetSuccessRate() float64 {
	total := m.GetLanguageSetTotal()
	if total == 0 {
		return 0.0
	}
	return float64(m.GetLanguageSetSuccess()) / float64(total)
}

// GetUploadSuccessRate returns the success rate for uploads (0.0 to 1.0).
func (m *Metrics) GetUploadSuccessRate() float64 {
	total := m.GetUploadTotal()
	if total == 0 {
		return 0.0
	}
	return float64(m.GetUploadSuccess()) / float64(total)
}

// Reset resets all metrics to zero.
func (m *Metrics) Reset() {
	atomic.StoreInt64(&m.LanguageSetSuccess, 0)
	atomic.StoreInt64(&m.LanguageSetFailure, 0)
	atomic.StoreInt64(&m.UploadSuccess, 0)
	atomic.StoreInt64(&m.UploadFailure, 0)
	atomic.StoreInt64(&m.LanguageValidation, 0)
	atomic.StoreInt64(&m.LanguageFallback, 0)
}
