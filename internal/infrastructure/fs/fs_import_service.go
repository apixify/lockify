package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/apixify/lockify/internal/domain/service"
)

type FsImportService struct {
}

func NewFsImportService() service.ImportService {
	return &FsImportService{}
}

func (service *FsImportService) FromJson(r io.Reader) (map[string]string, error) {
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

func (service *FsImportService) FromDotEnv(r io.Reader) (map[string]string, error) {
	entries := make(map[string]string)
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if len(value) >= 2 {
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
