package configuration

import (
	"os"
	"testing"
)

func TestConfig_WriteAndLoad(t *testing.T) {
	// Use a temporary file for testing
	tempFile := "./test_config.json"
	oldLocation := configLocation
	configLocation = tempFile
	defer func() {
		configLocation = oldLocation
		os.Remove(tempFile)
	}()

	cfg := &Config{
		CalendarUrl: "https://example.com/cal.ics",
		JiraEmail:   "test@example.com",
		JiraToken:   "secret_token",
	}

	// Test Write
	err := cfg.Write()
	if err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(tempFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Test LoadFromFile
	newCfg := &Config{}
	err = newCfg.LoadFromFile()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify data
	if newCfg.CalendarUrl != cfg.CalendarUrl {
		t.Errorf("Expected CalendarUrl %s, got %s", cfg.CalendarUrl, newCfg.CalendarUrl)
	}
	if newCfg.JiraEmail != cfg.JiraEmail {
		t.Errorf("Expected JiraEmail %s, got %s", cfg.JiraEmail, newCfg.JiraEmail)
	}
	if newCfg.JiraToken != cfg.JiraToken {
		t.Errorf("Expected JiraToken %s, got %s", cfg.JiraToken, newCfg.JiraToken)
	}
}

func TestConfig_LoadFromFile_NotFound(t *testing.T) {
	// Use a non-existent file
	tempFile := "./non_existent_config.json"
	oldLocation := configLocation
	configLocation = tempFile
	defer func() {
		configLocation = oldLocation
	}()

	cfg := &Config{}
	err := cfg.LoadFromFile()
	if err == nil {
		t.Error("Expected error when loading non-existent file, got nil")
	}
}
