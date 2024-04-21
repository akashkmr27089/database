package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Define a sample object structure
type Person struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Country string `json:"country"`
}

type Data struct {
	Data []Person `json:"data"`
}

func main2() {
	// Open the JSON file for appending
	file, err := os.OpenFile("persons.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer file.Close()
	var isEmptyFile bool = true

	// Check if the file is empty or not
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file information:", err)
		return
	}

	// If the file is not empty, add a comma before appending new data
	if fileInfo.Size() > 0 {
		isEmptyFile = false
		// if _, err := file.WriteString(","); err != nil {
		// 	fmt.Println("Error appending comma to JSON file:", err)
		// 	return
		// }
	}

	// Create a new Person object
	newPerson := Person{Name: "Charlie", Age: 35, Country: "Canada"}

	// Marshal the new Person object into JSON format
	newData, err := json.Marshal(newPerson)
	if err != nil {
		fmt.Println("Error marshalling new data to JSON:", err)
		return
	}

	// Append a comma and the new data to the JSON file
	if !isEmptyFile {
		if _, err := file.WriteString("," + string(newData)); err != nil {
			fmt.Println("Error appending data to JSON file:", err)
			return
		}
	} else {
		if _, err := file.WriteString(string(newData)); err != nil {
			fmt.Println("Error appending data to JSON file:", err)
			return
		}
	}

	fmt.Println("Data has been appended to persons.json")
}

// Could be later used for reading file and returning object
// Loading database to the inmemory
func readFileAndReturnObject(fileName string) {
	existingData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading existing JSON file:", err)
		return
	}

	if existingData != nil {
		val := fmt.Sprintf("{\"data\":["+"%s"+"]}", existingData)
		existingData2 := []byte(val)

		var existingPersons Data
		if err := json.Unmarshal(existingData2, &existingPersons); err != nil {
			fmt.Println("Error decoding existing JSON:", err)
			return
		}
	}
}
