package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestYAMLParsing(t *testing.T) {
	// Simple direct test of yaml library functionality
	yamlContent := []byte("name: Test Video\ncategory: testing\npath: /path/to/video.yaml\n")
	var video Video
	err := yaml.Unmarshal(yamlContent, &video)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	// Print the video struct to debug
	fmt.Printf("Parsed Video: %+v\n", video)

	if video.Name != "Test Video" {
		t.Errorf("Expected Name to be 'Test Video', got '%s'", video.Name)
	}
	if video.Category != "testing" {
		t.Errorf("Expected Category to be 'testing', got '%s'", video.Category)
	}

	// Try with struct literals to verify the Video struct is working
	directVideo := Video{
		Name:     "Test Video",
		Category: "testing",
	}

	if directVideo.Name != "Test Video" {
		t.Errorf("Direct assignment test failed. Expected Name to be 'Test Video', got '%s'", directVideo.Name)
	}
}

// TestExportedFieldParsing tests if the issue might be with lowercase vs uppercase field names
func TestExportedFieldParsing(t *testing.T) {
	// Test structure without explicit yaml tags - relying on yaml library's auto-conversion
	type TestVideo struct {
		Name     string
		Category string
		Path     string
	}

	yamlContent := []byte("name: Test Video\ncategory: testing\npath: /path/to/video.yaml\n")
	var video TestVideo
	err := yaml.Unmarshal(yamlContent, &video)
	if err != nil {
		t.Fatalf("Failed to unmarshal YAML: %v", err)
	}

	fmt.Printf("Parsed TestVideo: %+v\n", video)

	if video.Name != "Test Video" {
		t.Errorf("Expected Name to be 'Test Video', got '%s'", video.Name)
	}
}

// TestGetVideo tests the GetVideo functionality
func TestGetVideo(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test YAML file
	testPath := filepath.Join(tempDir, "test-video.yaml")
	testVideo := Video{
		Name:     "Test Video",
		Category: "testing",
		Path:     "/path/to/video.yaml",
	}

	// Write the YAML file
	y := YAML{}
	if err := y.WriteVideo(testVideo, testPath); err != nil {
		t.Fatalf("Failed to write test video YAML in TestGetVideo: %v", err)
	}

	// Read the YAML file
	video, err := y.GetVideo(testPath)
	if err != nil {
		t.Fatalf("GetVideo returned an error: %v", err)
	}

	// Verify the video was read correctly
	if video.Name != "Test Video" {
		t.Errorf("Expected video name to be 'Test Video', got '%s'", video.Name)
	}
	if video.Category != "testing" {
		t.Errorf("Expected video category to be 'testing', got '%s'", video.Category)
	}
	if video.Path != "/path/to/video.yaml" {
		t.Errorf("Expected video path to be '/path/to/video.yaml', got '%s'", video.Path)
	}
}

// TestWriteVideo tests the WriteVideo functionality
func TestWriteVideo(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "yaml-write-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test video
	testPath := filepath.Join(tempDir, "test-write-video.yaml")
	testVideo := Video{
		Name:     "Test Write Video",
		Category: "testing",
		Path:     "/path/to/written/video.yaml",
	}

	// Write the video to YAML
	y := YAML{}
	if err := y.WriteVideo(testVideo, testPath); err != nil {
		t.Fatalf("Failed to write test video YAML for TestWriteVideo: %v", err)
	}

	// Verify the file was created
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", testPath)
	}

	// Read the file back
	readVideo, err := y.GetVideo(testPath)
	if err != nil {
		t.Fatalf("GetVideo returned an error during read back: %v", err)
	}

	// Verify the contents
	if readVideo.Name != "Test Write Video" {
		t.Errorf("Expected video name to be 'Test Write Video', got '%s'", readVideo.Name)
	}
	if readVideo.Category != "testing" {
		t.Errorf("Expected video category to be 'testing', got '%s'", readVideo.Category)
	}
	if readVideo.Path != "/path/to/written/video.yaml" {
		t.Errorf("Expected video path to be '/path/to/written/video.yaml', got '%s'", readVideo.Path)
	}
}

// TestGetIndex tests the GetIndex functionality
func TestGetIndex(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "yaml-index-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test index file
	testPath := filepath.Join(tempDir, "index.json")

	// Create a simple index file
	indexContent := `[
		{"name": "Test Video 1", "category": "testing"},
		{"name": "Test Video 2", "category": "testing"}
	]`
	err = os.WriteFile(testPath, []byte(indexContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test index file: %v", err)
	}

	// Read the index
	y := YAML{
		IndexPath: testPath,
	}
	index, err := y.GetIndex()
	if err != nil {
		t.Fatalf("GetIndex returned an error: %v", err)
	}

	// Verify the index was read correctly
	if len(index) != 2 {
		t.Errorf("Expected index to have 2 entries, got %d", len(index))
	}
	if index[0].Name != "Test Video 1" {
		t.Errorf("Expected first video name to be 'Test Video 1', got '%s'", index[0].Name)
	}
	if index[1].Name != "Test Video 2" {
		t.Errorf("Expected second video name to be 'Test Video 2', got '%s'", index[1].Name)
	}
}

// TestWriteIndex tests the WriteIndex functionality
func TestWriteIndex(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "yaml-write-index-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test index
	testPath := filepath.Join(tempDir, "write-index.json")
	testIndex := []VideoIndex{
		{Name: "Test Write Video 1", Category: "testing"},
		{Name: "Test Write Video 2", Category: "testing"},
	}

	// Write the index
	y := YAML{
		IndexPath: testPath,
	}
	y.WriteIndex(testIndex)

	// Verify the file was created
	if _, err := os.Stat(testPath); os.IsNotExist(err) {
		t.Errorf("Expected file %s to exist, but it doesn't", testPath)
	}

	// Read the file back
	readIndex, err := y.GetIndex()
	if err != nil {
		t.Fatalf("GetIndex returned an error during read back: %v", err)
	}

	// Verify the contents
	if len(readIndex) != 2 {
		t.Errorf("Expected index to have 2 entries, got %d", len(readIndex))
	}
	if readIndex[0].Name != "Test Write Video 1" {
		t.Errorf("Expected first video name to be 'Test Write Video 1', got '%s'", readIndex[0].Name)
	}
	if readIndex[1].Name != "Test Write Video 2" {
		t.Errorf("Expected second video name to be 'Test Write Video 2', got '%s'", readIndex[1].Name)
	}
}

// TestNewYAML tests the NewYAML functionality
func TestNewYAML(t *testing.T) {
	// Create a YAML instance
	indexPath := "test-index.yaml"
	y := NewYAML(indexPath)

	// Verify it's not nil
	if y == nil {
		t.Errorf("Expected NewYAML to return a non-nil instance")
	}

	// Verify the index path is set correctly
	if y.IndexPath != indexPath {
		t.Errorf("Expected IndexPath to be '%s', got '%s'", indexPath, y.IndexPath)
	}

	// Test that NewYAML creates a YAML struct with the correct IndexPath
	testIndexPath := "/test/path/index.json"
	newY := NewYAML(testIndexPath)
	if newY.IndexPath != testIndexPath {
		t.Errorf("Expected IndexPath to be '%s', got '%s'", testIndexPath, newY.IndexPath)
	}
}

func TestGetVideo_FileNotFound(t *testing.T) {
	y := YAML{}
	_, err := y.GetVideo("non_existent_path.yaml")
	if err == nil {
		t.Fatalf("Expected GetVideo to return an error for non-existent file, but got nil")
	}
	// Check if the error is an os.PathError, which is what os.ReadFile returns for non-existent files
	if !os.IsNotExist(err) {
		// It might be wrapped, so check unwrap
		type unwrap interface {
			Unwrap() error
		}
		if unwrapErr, ok := err.(unwrap); ok {
			if !os.IsNotExist(unwrapErr.Unwrap()) {
				t.Errorf("Expected GetVideo to return an os.IsNotExist error, got %T: %v", err, err)
			}
		} else {
			t.Errorf("Expected GetVideo to return an os.IsNotExist error, got %T: %v", err, err)
		}
	}
}

func TestGetVideo_InvalidYAML(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "invalid-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	invalidYAMLPath := filepath.Join(tempDir, "invalid.yaml")
	if err := os.WriteFile(invalidYAMLPath, []byte("name: Test Video\ncategory: testing\n  badlyIndentedKey: true"), 0644); err != nil {
		t.Fatalf("Failed to write invalid YAML file: %v", err)
	}

	y := YAML{}
	_, err = y.GetVideo(invalidYAMLPath)
	if err == nil {
		t.Fatalf("Expected GetVideo to return an error for invalid YAML, but got nil")
	}
	// We expect an error from yaml.Unmarshal, check for it.
	// The error message from our function is "failed to unmarshal video data from %s: %w"
	expectedErrorMsgPart := "failed to unmarshal video data"
	if !strings.Contains(err.Error(), expectedErrorMsgPart) {
		t.Errorf("Expected GetVideo error to contain '%s', got '%s'", expectedErrorMsgPart, err.Error())
	}
}

func TestGetIndex_FileNotFound(t *testing.T) {
	y := YAML{IndexPath: "non_existent_index.json"}
	_, err := y.GetIndex()
	if err == nil {
		t.Fatalf("Expected GetIndex to return an error for non-existent file, but got nil")
	}
	if !os.IsNotExist(err) {
		// It might be wrapped, so check unwrap
		type unwrap interface {
			Unwrap() error
		}
		if unwrapErr, ok := err.(unwrap); ok {
			if !os.IsNotExist(unwrapErr.Unwrap()) {
				t.Errorf("Expected GetIndex to return an os.IsNotExist error, got %T: %v", err, err)
			}
		} else {
			t.Errorf("Expected GetIndex to return an os.IsNotExist error, got %T: %v", err, err)
		}
	}
}

func TestGetIndex_InvalidYAML(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "invalid-index-yaml-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	invalidIndexYAMLPath := filepath.Join(tempDir, "invalid_index.yaml")
	if err := os.WriteFile(invalidIndexYAMLPath, []byte("[{\"name\": \"Test Video 1\", \"category\": \"testing\"}, {invalid_json]"), 0644); err != nil {
		t.Fatalf("Failed to write invalid index YAML file: %v", err)
	}

	y := YAML{IndexPath: invalidIndexYAMLPath}
	_, err = y.GetIndex()
	if err == nil {
		t.Fatalf("Expected GetIndex to return an error for invalid YAML, but got nil")
	}
	expectedErrorMsgPart := "failed to unmarshal video index"
	if !strings.Contains(err.Error(), expectedErrorMsgPart) {
		t.Errorf("Expected GetIndex error to contain '%s', got '%s'", expectedErrorMsgPart, err.Error())
	}
}

func TestVideo_JSONConsistency(t *testing.T) {
	t.Run("Video struct should serialize to camelCase JSON", func(t *testing.T) {
		video := Video{
			Name:        "test-video",
			ProjectName: "Test Project",
			ProjectURL:  "https://example.com",
			Sponsorship: Sponsorship{
				Amount:  "1000",
				Emails:  "sponsor@example.com",
				Blocked: "false",
			},
		}

		// Test serialization (GET response behavior)
		jsonData, err := json.Marshal(video)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(jsonData, &jsonMap)
		require.NoError(t, err)

		// Should be camelCase, not PascalCase
		assert.Equal(t, "Test Project", jsonMap["projectName"])
		assert.Equal(t, "https://example.com", jsonMap["projectURL"])

		// Sponsorship nested fields should also be camelCase
		sponsorship, ok := jsonMap["sponsorship"].(map[string]interface{})
		require.True(t, ok, "sponsorship should be a JSON object")
		assert.Equal(t, "1000", sponsorship["amount"])
		assert.Equal(t, "sponsor@example.com", sponsorship["emails"])
		assert.Equal(t, "false", sponsorship["blocked"])

		// These PascalCase fields should NOT exist
		assert.NotContains(t, jsonMap, "ProjectName")
		assert.NotContains(t, jsonMap, "ProjectURL")
	})

	t.Run("Video struct should deserialize from camelCase JSON", func(t *testing.T) {
		// Test deserialization (PUT request behavior)
		jsonData := `{
			"name": "test-video",
			"projectName": "Test Project",
			"projectURL": "https://example.com",
			"sponsorship": {
				"amount": "1000",
				"emails": "sponsor@example.com",
				"blocked": "false"
			}
		}`

		var video Video
		err := json.Unmarshal([]byte(jsonData), &video)
		require.NoError(t, err)

		assert.Equal(t, "test-video", video.Name)
		assert.Equal(t, "Test Project", video.ProjectName)
		assert.Equal(t, "https://example.com", video.ProjectURL)
		assert.Equal(t, "1000", video.Sponsorship.Amount)
		assert.Equal(t, "sponsor@example.com", video.Sponsorship.Emails)
		assert.Equal(t, "false", video.Sponsorship.Blocked)
	})

}

// TestVideo_GetLanguage tests the GetLanguage helper method
func TestVideo_GetLanguage(t *testing.T) {
	tests := []struct {
		name           string
		video          Video
		defaultLang    string
		expectedResult string
	}{
		{
			name: "Language field set, should return video language",
			video: Video{
				Name:     "Test Video",
				Language: "es",
			},
			defaultLang:    "en",
			expectedResult: "es",
		},
		{
			name: "Language field empty, should return default",
			video: Video{
				Name:     "Test Video",
				Language: "",
			},
			defaultLang:    "fr",
			expectedResult: "fr",
		},
		{
			name: "Language field empty, default empty, should return empty",
			video: Video{
				Name:     "Test Video",
				Language: "",
			},
			defaultLang:    "",
			expectedResult: "",
		},
		{
			name: "Language field set to same as default",
			video: Video{
				Name:     "Test Video",
				Language: "de",
			},
			defaultLang:    "de",
			expectedResult: "de",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.video.GetLanguage(tt.defaultLang)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// TestVideo_GetAudioLanguage tests the GetAudioLanguage helper method
func TestVideo_GetAudioLanguage(t *testing.T) {
	tests := []struct {
		name           string
		video          Video
		defaultLang    string
		expectedResult string
	}{
		{
			name: "AudioLanguage field set, should return video audio language",
			video: Video{
				Name:          "Test Video",
				AudioLanguage: "ja",
			},
			defaultLang:    "en",
			expectedResult: "ja",
		},
		{
			name: "AudioLanguage field empty, should return default",
			video: Video{
				Name:          "Test Video",
				AudioLanguage: "",
			},
			defaultLang:    "ko",
			expectedResult: "ko",
		},
		{
			name: "AudioLanguage field empty, default empty, should return empty",
			video: Video{
				Name:          "Test Video",
				AudioLanguage: "",
			},
			defaultLang:    "",
			expectedResult: "",
		},
		{
			name: "AudioLanguage field set to same as default",
			video: Video{
				Name:          "Test Video",
				AudioLanguage: "zh",
			},
			defaultLang:    "zh",
			expectedResult: "zh",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.video.GetAudioLanguage(tt.defaultLang)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

// TestVideo_LanguageSerialization tests serialization/deserialization with language fields
func TestVideo_LanguageSerialization(t *testing.T) {
	t.Run("Video with language fields should serialize correctly", func(t *testing.T) {
		video := Video{
			Name:          "Test Video",
			Language:      "es",
			AudioLanguage: "fr",
		}

		// Test JSON serialization
		jsonData, err := json.Marshal(video)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(jsonData, &jsonMap)
		require.NoError(t, err)

		assert.Equal(t, "es", jsonMap["language"])
		assert.Equal(t, "fr", jsonMap["audioLanguage"])
	})

	t.Run("Video without language fields should serialize with empty strings", func(t *testing.T) {
		video := Video{
			Name: "Test Video",
			// Language and AudioLanguage are empty strings by default
		}

		jsonData, err := json.Marshal(video)
		require.NoError(t, err)

		var jsonMap map[string]interface{}
		err = json.Unmarshal(jsonData, &jsonMap)
		require.NoError(t, err)

		// Check if language fields exist and are empty strings
		// Note: JSON serialization might omit empty strings depending on tags
		language, hasLanguage := jsonMap["language"]
		audioLanguage, hasAudioLanguage := jsonMap["audioLanguage"]

		if hasLanguage {
			assert.Equal(t, "", language)
		}
		if hasAudioLanguage {
			assert.Equal(t, "", audioLanguage)
		}
	})

	t.Run("Video should deserialize from JSON with language fields", func(t *testing.T) {
		jsonData := `{
			"name": "Test Video",
			"language": "de",
			"audioLanguage": "it"
		}`

		var video Video
		err := json.Unmarshal([]byte(jsonData), &video)
		require.NoError(t, err)

		assert.Equal(t, "Test Video", video.Name)
		assert.Equal(t, "de", video.Language)
		assert.Equal(t, "it", video.AudioLanguage)
	})
}

// TestVideo_LanguageYAMLSerialization tests YAML serialization/deserialization with language fields
func TestVideo_LanguageYAMLSerialization(t *testing.T) {
	t.Run("Video with language fields should serialize to YAML correctly", func(t *testing.T) {
		video := Video{
			Name:          "Test Video",
			Language:      "pt",
			AudioLanguage: "ru",
		}

		// Test YAML serialization
		yamlData, err := yaml.Marshal(video)
		require.NoError(t, err)

		// Verify YAML contains language fields
		yamlStr := string(yamlData)
		assert.Contains(t, yamlStr, "language: pt")
		assert.Contains(t, yamlStr, "audioLanguage: ru")
	})

	t.Run("Video should deserialize from YAML with language fields", func(t *testing.T) {
		yamlData := `name: Test Video
language: nl
audioLanguage: sv`

		var video Video
		err := yaml.Unmarshal([]byte(yamlData), &video)
		require.NoError(t, err)

		assert.Equal(t, "Test Video", video.Name)
		assert.Equal(t, "nl", video.Language)
		assert.Equal(t, "sv", video.AudioLanguage)
	})
}

// TestVideo_BackwardCompatibility tests backward compatibility with existing metadata
func TestVideo_BackwardCompatibility(t *testing.T) {
	t.Run("Existing video without language fields should work with new methods", func(t *testing.T) {
		// Simulate an old video metadata without language fields
		oldVideo := Video{
			Name:     "Old Video",
			Category: "testing",
			Title:    "Old Title",
			// Language and AudioLanguage are empty strings (default)
		}

		// Should not panic and should return defaults
		language := oldVideo.GetLanguage("en")
		audioLanguage := oldVideo.GetAudioLanguage("en")

		assert.Equal(t, "en", language)
		assert.Equal(t, "en", audioLanguage)
	})

	t.Run("Video with partial language fields should work correctly", func(t *testing.T) {
		video := Video{
			Name:     "Partial Language Video",
			Language: "es",
			// AudioLanguage is empty
		}

		// Language should return video's language
		language := video.GetLanguage("en")
		assert.Equal(t, "es", language)

		// AudioLanguage should return default
		audioLanguage := video.GetAudioLanguage("fr")
		assert.Equal(t, "fr", audioLanguage)
	})
}
