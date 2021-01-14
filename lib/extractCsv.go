package lib

import (
	"encoding/csv"
	"os"
)

// ExtractCsv
func ExtractCsv(affiliate, fileName string, maps map[string]string) bool {
	file, err := os.Create(fileName)
	if err != nil {
		return false
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{affiliate, ""}
	err = writer.Write(headers)

	for key, value := range maps {
		err = writer.Write([]string{key, value})
	}

	return true
}
