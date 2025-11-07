package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gosimple/slug"
)

func FormatLocaleResource(inputData map[string]string) map[string]map[string]interface{} {
	outputData := map[string]map[string]interface{}{}

	for key, value := range inputData {
		prefix := strings.ToLower(value)
		prefix = slug.Make(prefix)
		prefix = strings.ReplaceAll(prefix, "-", " ")

		outputData[value] = map[string]interface{}{
			"prefix": prefix,
			"body":   []string{key},
		}
	}

	return outputData
}

func WriteResourceToSnippetFile(data map[string]map[string]interface{}, filePath string) error {
	dir := filepath.Dir(filePath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create directory for resource snippet: %v", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Failed to create file for resource snippet: %v", err)
	}
	defer file.Close()

	snippetData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return fmt.Errorf("Failed to marshal JSON: %v", err)
	}

	_, err = file.Write(snippetData)
	if err != nil {
		return fmt.Errorf("Failed to write to file: %v", err)
	}

	return nil
}
