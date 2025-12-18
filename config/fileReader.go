package config

import (
	"encoding/csv"
	"os"

	"github.com/gofiber/fiber/v2"
)

// FileExists checks if a file exists at the given path
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// ReadCSV reads a CSV file and returns records as a slice of string slices
func ReadCSV(c *fiber.Ctx, filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		LogError(c, "Failed to open CSV file", "file_path", filePath, "error", err.Error())
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		LogError(c, "Failed to read CSV file", "file_path", filePath, "error", err.Error())
		return nil, err
	}

	LogInfo(c, "CSV file read successfully", "file_path", filePath, "records", len(records))
	return records, nil
}

// ReadCSVWithHeaders reads a CSV file and returns records as a slice of maps with headers as keys
func ReadCSVWithHeaders(c *fiber.Ctx, filePath string) ([]map[string]string, error) {
	records, err := ReadCSV(c, filePath)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []map[string]string{}, nil
	}

	headers := records[0]
	var result []map[string]string

	for _, record := range records[1:] {
		row := make(map[string]string)
		for i, header := range headers {
			if i < len(record) {
				row[header] = record[i]
			}
		}
		result = append(result, row)
	}

	return result, nil
}
