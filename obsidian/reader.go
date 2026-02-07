package obsidian

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetDailyNotePath returns the expected path for a daily note based on the date.
// Format: {{YEAR}}/YYYY-MM-DD.md
func GetDailyNotePath(vaultPath string, date time.Time) string {
	year := date.Format("2006")
	filename := date.Format("2006-01-02") + ".md"
	return filepath.Join(vaultPath, year, filename)
}

// GetWeeklyNotePath returns the expected path for a weekly note based on the date.
// Format: {{YEAR}}/Week WW.md
func GetWeeklyNotePath(vaultPath string, date time.Time) string {
	year := date.Format("2006")
	_, week := date.ISOWeek()
	filename := fmt.Sprintf("Week %02d.md", week)
	return filepath.Join(vaultPath, year, filename)
}

// ReadSection reads a specific section from a markdown file starting from the given header.
// It captures everything under the header until the next header of the same level (##).
func ReadSection(filePath string, header string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var content strings.Builder
	foundHeader := false

	// Target header should look like "## Header"
	targetHeader := "## " + header

	for scanner.Scan() {
		line := scanner.Text()

		if !foundHeader {
			if strings.TrimSpace(line) == targetHeader {
				foundHeader = true
				content.WriteString(line + "\n")
			}
			continue
		}

		// If we found the header, we start capturing until we see another level 2 header
		if strings.HasPrefix(line, "## ") {
			break
		}

		content.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	if !foundHeader {
		return "", fmt.Errorf("header %q not found in %s", header, filePath)
	}

	return content.String(), nil
}
