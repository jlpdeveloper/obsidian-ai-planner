package obsidian

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetDailyNotePath(t *testing.T) {
	vault := "C:\\Vault"
	date := time.Date(2026, 2, 7, 0, 0, 0, 0, time.UTC)
	expected := filepath.Join(vault, "2026", "2026-02-07.md")
	result := GetDailyNotePath(vault, date)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestGetWeeklyNotePath(t *testing.T) {
	vault := "C:\\Vault"
	date := time.Date(2026, 2, 7, 0, 0, 0, 0, time.UTC)
	// Feb 7, 2026 is Saturday of Week 06 (ISO)
	expected := filepath.Join(vault, "2026", "Week 06.md")
	result := GetWeeklyNotePath(vault, date)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestReadSection(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.md")
	content := `---
journal-date: 2026-02-07
---

## Tasks
### Goals
- [ ] Task 1
### Meetings
- [ ] Meeting 1

## Other Section
Some other content
`
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	section, err := ReadSection(tmpFile, "Tasks")
	if err != nil {
		t.Fatalf("ReadSection failed: %v", err)
	}

	expectedSection := "## Tasks\n### Goals\n- [ ] Task 1\n### Meetings\n- [ ] Meeting 1\n\n"
	if section != expectedSection {
		t.Errorf("Expected:\n%q\nGot:\n%q", expectedSection, section)
	}
}

func TestReadSection_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_not_found.md")
	content := "## Something Else\nContent\n"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ReadSection(tmpFile, "Tasks")
	if err == nil {
		t.Error("Expected error for missing section, got nil")
	}
}
