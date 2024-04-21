package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func getFileName(uuid string) string {
	filename := fmt.Sprintf("%s.json", uuid)
	return filename
}

// Based on uuid, append all the operation files to a single file
func AppendOperationsIntoFileBasedOnUuid(
	uuid string,
	operations []Operation,
) (bool, error) {
	fmt.Printf("%s-%v", uuid, operations)
	if len(operations) == 0 {
		return false, nil
	}

	filename := getFileName(uuid)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return false, err
	}
	defer file.Close()

	var isEmptyFile bool = true

	// Check if the file is empty or not
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return false, err
	}

	if fileInfo.Size() > 0 {
		isEmptyFile = false
	}

	arrayOfOperationString := make([]string, len(operations))
	for idx, operation := range operations {
		newData, err := json.Marshal(operation)
		if err != nil {
			fmt.Println("Error marshalling new data to JSON:", err)
			return false, err
		}

		arrayOfOperationString[idx] = string(newData)
	}

	finalString := strings.Join(arrayOfOperationString, ",")

	// Append a comma and the new data to the JSON file
	if !isEmptyFile {
		if _, err := file.WriteString("," + string(finalString)); err != nil {
			fmt.Println("Error appending data to JSON file:", err)
			return false, err
		}
	} else {
		if _, err := file.WriteString(string(finalString)); err != nil {
			fmt.Println("Error appending data to JSON file:", err)
			return false, err
		}
	}

	return true, nil
}
