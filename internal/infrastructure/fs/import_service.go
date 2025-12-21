package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
)

const (
	// dotenvKeyValueParts is the expected number of parts when splitting a dotenv line by "=".
	dotenvKeyValueParts = 2
	// minValueLength is the minimum length for a value to be considered (needs at least 1 char).
	minValueLength = 2
)

// ImportService implements ImportService for filesystem-based imports.
type ImportService struct{}

// NewImportService creates a new ImportService instance.
func NewImportService() service.ImportService {
	return &ImportService{}
}

// FromJSON parses JSON data from a reader and returns a map of key-value pairs.
func (service *ImportService) FromJSON(r io.Reader) (map[string]string, error) {
	var entries map[string]string
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&entries); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	result := make(map[string]string, len(entries))
	for k, v := range entries {
		result[k] = fmt.Sprintf("%v", v)
	}

	return result, nil
}

// FromDotEnv parses dotenv data from a reader and returns a map of key-value pairs.
func (service *ImportService) FromDotEnv(r io.Reader) (map[string]string, error) {
	entries := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.SplitN(line, "=", dotenvKeyValueParts)
		if len(parts) != dotenvKeyValueParts {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if len(value) >= minValueLength {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		if key != "" {
			entries[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return entries, nil
}
