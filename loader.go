package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type OperationData struct {
	Data []Operation `json:"data"`
}

// Function -> to dump in the collection

// Could be later used for reading file and returning object
// Loading database to the inmemory
func ReadFileAndReturnObject(uuid string) (*OperationData, error) {
	filename := getFileName(uuid)
	existingData, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading existing JSON file:", err)
		return nil, err
	}

	var existingOperationData OperationData
	if existingData != nil {
		val := fmt.Sprintf("{\"data\":["+"%s"+"]}", existingData)
		existingData2 := []byte(val)

		if err := json.Unmarshal(existingData2, &existingOperationData); err != nil {
			fmt.Println("Error decoding existing JSON:", err)
			return nil, err
		}
	}

	return &existingOperationData, nil
}
