package publishing

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetrics_Counters(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	// Test initial state
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageSetSuccess())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageSetFailure())
	assert.Equal(t, int64(0), YouTubeMetrics.GetUploadSuccess())
	assert.Equal(t, int64(0), YouTubeMetrics.GetUploadFailure())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageValidation())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageFallback())

	// Test incrementing counters
	YouTubeMetrics.IncLanguageSetSuccess()
	YouTubeMetrics.IncLanguageSetFailure()
	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadFailure()
	YouTubeMetrics.IncLanguageValidation()
	YouTubeMetrics.IncLanguageFallback()

	// Test counter values
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageSetSuccess())
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageSetFailure())
	assert.Equal(t, int64(1), YouTubeMetrics.GetUploadSuccess())
	assert.Equal(t, int64(1), YouTubeMetrics.GetUploadFailure())
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageValidation())
	assert.Equal(t, int64(1), YouTubeMetrics.GetLanguageFallback())
}

func TestMetrics_ConcurrentAccess(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	// Test concurrent access
	const numGoroutines = 100
	const incrementsPerGoroutine = 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				YouTubeMetrics.IncLanguageSetSuccess()
				YouTubeMetrics.IncUploadSuccess()
			}
		}()
	}

	wg.Wait()

	// Verify final counts
	expectedCount := int64(numGoroutines * incrementsPerGoroutine)
	assert.Equal(t, expectedCount, YouTubeMetrics.GetLanguageSetSuccess())
	assert.Equal(t, expectedCount, YouTubeMetrics.GetUploadSuccess())
}

func TestMetrics_TotalCalculations(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	// Set up test data
	YouTubeMetrics.IncLanguageSetSuccess()
	YouTubeMetrics.IncLanguageSetSuccess()
	YouTubeMetrics.IncLanguageSetFailure()

	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadFailure()

	// Test total calculations
	assert.Equal(t, int64(3), YouTubeMetrics.GetLanguageSetTotal())
	assert.Equal(t, int64(4), YouTubeMetrics.GetUploadTotal())
}

func TestMetrics_SuccessRates(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	// Test with zero values
	assert.Equal(t, 0.0, YouTubeMetrics.GetLanguageSetSuccessRate())
	assert.Equal(t, 0.0, YouTubeMetrics.GetUploadSuccessRate())

	// Test with some successes and failures
	YouTubeMetrics.IncLanguageSetSuccess()
	YouTubeMetrics.IncLanguageSetSuccess()
	YouTubeMetrics.IncLanguageSetFailure()

	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadFailure()

	// Test success rates
	assert.Equal(t, 2.0/3.0, YouTubeMetrics.GetLanguageSetSuccessRate())
	assert.Equal(t, 3.0/4.0, YouTubeMetrics.GetUploadSuccessRate())
}

func TestMetrics_Reset(t *testing.T) {
	// Set up some data
	YouTubeMetrics.IncLanguageSetSuccess()
	YouTubeMetrics.IncLanguageSetFailure()
	YouTubeMetrics.IncUploadSuccess()
	YouTubeMetrics.IncUploadFailure()
	YouTubeMetrics.IncLanguageValidation()
	YouTubeMetrics.IncLanguageFallback()

	// Verify data exists
	assert.Greater(t, YouTubeMetrics.GetLanguageSetSuccess(), int64(0))
	assert.Greater(t, YouTubeMetrics.GetLanguageSetFailure(), int64(0))
	assert.Greater(t, YouTubeMetrics.GetUploadSuccess(), int64(0))
	assert.Greater(t, YouTubeMetrics.GetUploadFailure(), int64(0))
	assert.Greater(t, YouTubeMetrics.GetLanguageValidation(), int64(0))
	assert.Greater(t, YouTubeMetrics.GetLanguageFallback(), int64(0))

	// Reset metrics
	YouTubeMetrics.Reset()

	// Verify all counters are zero
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageSetSuccess())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageSetFailure())
	assert.Equal(t, int64(0), YouTubeMetrics.GetUploadSuccess())
	assert.Equal(t, int64(0), YouTubeMetrics.GetUploadFailure())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageValidation())
	assert.Equal(t, int64(0), YouTubeMetrics.GetLanguageFallback())
}

func TestMetrics_EdgeCases(t *testing.T) {
	// Reset metrics to ensure clean state
	YouTubeMetrics.Reset()

	// Test success rate with only failures
	YouTubeMetrics.IncLanguageSetFailure()
	YouTubeMetrics.IncLanguageSetFailure()

	assert.Equal(t, 0.0, YouTubeMetrics.GetLanguageSetSuccessRate())

	// Test success rate with only successes
	YouTubeMetrics.Reset()
	YouTubeMetrics.IncLanguageSetSuccess()
	YouTubeMetrics.IncLanguageSetSuccess()

	assert.Equal(t, 1.0, YouTubeMetrics.GetLanguageSetSuccessRate())
}
