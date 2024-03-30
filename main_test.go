package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_ENV", "test_value")
	defer os.Unsetenv("TEST_ENV")

	tests := []struct {
		name     string
		key      string
		fallback string
		want     string
	}{
		{"Existing environment variable", "TEST_ENV", "fallback", "test_value"},
		{"Non-existing environment variable", "NON_EXISTING_ENV", "fallback", "fallback"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getEnv(tt.key, tt.fallback); got != tt.want {
				t.Errorf("getEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Write a temporary YAML file for testing
	tmpfile, err := ioutil.TempFile("", "test.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte(`
sites:
  - name: Test Site
    url: http://test.com
`)
	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	config, err := loadConfig(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(config.Sites) != 1 {
		t.Errorf("Expected 1 site, got %d", len(config.Sites))
	}

	if config.Sites[0].Name != "Test Site" {
		t.Errorf("Expected site name 'Test Site', got '%s'", config.Sites[0].Name)
	}

	if config.Sites[0].URL != "http://test.com" {
		t.Errorf("Expected site URL 'http://test.com', got '%s'", config.Sites[0].URL)
	}
}

// Note: Testing isSiteUp function would require mocking the http client which is beyond the scope of this example.
// Similarly, testing writeMarkdown function would require setting up a proper template file and checking the output file content.
