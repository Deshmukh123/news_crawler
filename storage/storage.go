package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

// SaveToJSON writes the extracted articles into a JSON file
func SaveToJSON(filename string, data interface{}) error {
	// Create file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode JSON data
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	fmt.Println("📁 Data saved to", filename)
	return nil
}
