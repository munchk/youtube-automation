package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguageConstants(t *testing.T) {
	// Test that language constants are correctly defined
	assert.Equal(t, "en", DefaultLanguage, "DefaultLanguage should be 'en'")
	assert.Equal(t, "en", LanguageEnglish, "LanguageEnglish should be 'en'")
	assert.Equal(t, DefaultLanguage, LanguageEnglish, "DefaultLanguage and LanguageEnglish should be the same")
}

func TestLanguageMap(t *testing.T) {
	// Test that LanguageMap contains expected mappings
	assert.Equal(t, "English", LanguageMap[LanguageEnglish], "LanguageMap should map 'en' to 'English'")

	// Test that all language constants have corresponding entries in LanguageMap
	assert.Contains(t, LanguageMap, LanguageEnglish, "LanguageMap should contain LanguageEnglish key")

	// Test that LanguageMap is not empty
	assert.NotEmpty(t, LanguageMap, "LanguageMap should not be empty")
}

func TestIsValidLanguage(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		expected bool
	}{
		{
			name:     "Valid English code",
			code:     "en",
			expected: true,
		},
		{
			name:     "Valid English code from constant",
			code:     LanguageEnglish,
			expected: true,
		},
		{
			name:     "Valid English code from default constant",
			code:     DefaultLanguage,
			expected: true,
		},
		{
			name:     "Invalid language code",
			code:     "invalid",
			expected: false,
		},
		{
			name:     "Empty string",
			code:     "",
			expected: false,
		},
		{
			name:     "Spanish code (not supported yet)",
			code:     "es",
			expected: false,
		},
		{
			name:     "French code (not supported yet)",
			code:     "fr",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidLanguage(tt.code)
			assert.Equal(t, tt.expected, result, "IsValidLanguage(%s) should return %v", tt.code, tt.expected)
		})
	}
}

func TestLanguageConstantsConsistency(t *testing.T) {
	// Test that all constants are consistent
	assert.Equal(t, DefaultLanguage, LanguageEnglish, "DefaultLanguage and LanguageEnglish should be consistent")

	// Test that constants match LanguageMap keys
	assert.Contains(t, LanguageMap, DefaultLanguage, "LanguageMap should contain DefaultLanguage key")
	assert.Contains(t, LanguageMap, LanguageEnglish, "LanguageMap should contain LanguageEnglish key")

	// Test that the default language is valid
	assert.True(t, IsValidLanguage(DefaultLanguage), "DefaultLanguage should be valid")
	assert.True(t, IsValidLanguage(LanguageEnglish), "LanguageEnglish should be valid")
}

func TestLanguageMapStructure(t *testing.T) {
	// Test that LanguageMap has the expected structure
	expectedKeys := []string{LanguageEnglish}
	expectedValues := []string{"English"}

	for _, key := range expectedKeys {
		assert.Contains(t, LanguageMap, key, "LanguageMap should contain key: %s", key)
	}

	// Test that expected values exist in the map
	for _, expectedValue := range expectedValues {
		found := false
		for _, value := range LanguageMap {
			if value == expectedValue {
				found = true
				break
			}
		}
		assert.True(t, found, "LanguageMap should contain value: %s", expectedValue)
	}

	// Test that all values in LanguageMap are non-empty
	for key, value := range LanguageMap {
		assert.NotEmpty(t, key, "LanguageMap key should not be empty")
		assert.NotEmpty(t, value, "LanguageMap value for key '%s' should not be empty", key)
	}
}
