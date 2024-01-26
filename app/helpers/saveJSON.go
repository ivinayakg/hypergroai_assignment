package helpers

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

const filename = "client_service.json"

func SaveAsJSON(input string) {
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		if err != nil {
			log.Fatalf("Failed to delete the file. %v", err)
		}
	}

	// Convert string to JSON
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(input), &jsonData)
	if err != nil {
		log.Fatalf("Error occurred during unmarshalling. %v", err)
	}

	privateKey, ok := jsonData["private_key"].(string)
	if ok {
		jsonData["private_key"] = strings.Replace(privateKey, "%n%", "\n", -1)
	}

	// Convert back to JSON to format it
	formattedJson, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		log.Fatalf("Error occurred during marshalling. %v", err)
	}

	// Write to file
	err = os.WriteFile(filename, formattedJson, 0644)
	if err != nil {
		log.Fatalf("Error occurred during writing the file. %v", err)
	}
}
