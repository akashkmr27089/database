package main

import (
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
)

var RESET_COUNTER_CONSTANT int8 = 5

type Operation struct {
	Name  OperationName
	Key   *string
	Value *string
}

type Collection struct {
	name           string
	uuid           string
	updateStrategy UpdateStrategy
	Documents      map[string]string
	keys           []string
	Operations     []Operation
	flushCounter   int8
}

func (c *Collection) incrementFlashCounter(operation Operation) {
	if c.updateStrategy == BatchUpdateStrategy {
		c.flushCounter += 1
		if operation.Name == CreateCollectionOperationName {
			AppendOperationsIntoFileBasedOnUuid(c.uuid, []Operation{operation})
			return
		}
		if c.flushCounter%RESET_COUNTER_CONSTANT == 0 {
			// This will send all the new updates to Database
			c.Operations = append(c.Operations, operation)
			AppendOperationsIntoFileBasedOnUuid(c.uuid, c.Operations)

			// Reset all the update variables
			c.Operations = []Operation{}
			c.flushCounter = 0
		} else {
			c.Operations = append(c.Operations, operation)
		}
	} else if c.updateStrategy == IncrementalUpdateStrategy {
		AppendOperationsIntoFileBasedOnUuid(c.uuid, []Operation{operation})
	}
}

func (c *Collection) Create(name string, updateStrategy *UpdateStrategy, preUuid *string) {
	if updateStrategy == nil {
		c.updateStrategy = BatchUpdateStrategy
	}
	c.name = name
	if preUuid != nil {
		c.uuid = *preUuid
	} else {
		c.uuid = uuid.New().String()
	}
	c.Documents = make(map[string]string, 0)

	// Add Element in the Operation Array
	operation := Operation{
		Name:  CreateCollectionOperationName,
		Value: &name,
	}
	c.incrementFlashCounter(operation)
}

func (c *Collection) Close() {
	AppendOperationsIntoFileBasedOnUuid(c.uuid, c.Operations)
}

// Note that the validation for key and value must be done beforehand
func (c *Collection) InsertOne(key string, value string) bool {
	_, ok := c.Documents[key]
	var isInserted bool = false
	if !ok {
		c.Documents[key] = value
		c.clearkeyList()
		isInserted = true

		// Add Element in the Operation Array
		operation := Operation{
			Name:  InsertOneOperationName,
			Key:   &key,
			Value: &value,
		}
		c.incrementFlashCounter(operation)
	}

	return isInserted
}

func (c *Collection) UpdateOne(key string, value string) bool {
	_, ok := c.Documents[key]
	var isUpdated bool = false
	if ok {
		c.Documents[key] = value
		c.clearkeyList()
		isUpdated = true

		// Add Element in the Operation Array
		operation := Operation{
			Name:  UpdateOneOperationName,
			Key:   &key,
			Value: &value,
		}
		c.incrementFlashCounter(operation)
	}

	return isUpdated
}

func (c *Collection) DeleteOne(key string) bool {
	_, ok := c.Documents[key]
	var isDeleted bool = false
	if !ok {
		delete(c.Documents, key)
		c.clearkeyList()
		isDeleted = true

		// Add Element in the Operation Array
		operation := Operation{
			Name: DeleteOneOperationName,
			Key:  &key,
		}
		c.incrementFlashCounter(operation)
	}

	return isDeleted
}

// todo: need to update function clearkeyList, populateKeyList, getAllKeys
// Needed this function to clear the keys list whenever update, create, delete operation
func (c *Collection) clearkeyList() {
	if len(c.keys) > 0 {
		c.keys = []string{}
	}
}

func (c *Collection) populateKeyList() {
	for k := range c.Documents {
		c.keys = append(c.keys, k)
	}
}

func (c *Collection) getAllKeys() []string {
	if len(c.keys) == 0 {
		c.populateKeyList()
	}
	return c.keys
}

func (c *Collection) GetAllKeys() []string {
	keys := c.getAllKeys()
	return keys
}

// Requirement for storing the keys value
// Should be able to run regex in it
// Should be able to remove the variable easily -->

// TODO: Temporary implementaion
// Get keys by regex
func (c *Collection) GetKeysByParitialByRegex(regexKey string) []string {
	keys := c.getAllKeys()

	filteredKeys := make([]string, 0)
	for _, str := range keys {
		matched, err := filepath.Match(regexKey, str)
		if err != nil {
			// Handle error
			fmt.Println("Error:", err)
			continue
		}
		if matched {
			filteredKeys = append(filteredKeys, str)
		}
	}

	return filteredKeys
}

func main() {
	// If the file exists --> Load the Data to the collections
	var c Collection
	// loadPrexistingData()
	uuid := "4a13d3ab-ebda-4d76-aebb-b7a0f7ce2bbb"
	c.Create("Test", nil, &uuid)

	c.InsertOne("test", "val")
	c.InsertOne("test1", "val")
	c.InsertOne("test12", "val2")
	c.InsertOne("test3", "val34")
	c.InsertOne("test3", "val")
	c.InsertOne("test13", "val")
	c.InsertOne("test123", "val2")
	c.InsertOne("test33", "val34")

	values := c.GetKeysByParitialByRegex("*test1*")
	for _, value := range values {
		fmt.Println(value)
	}

	c.Close()

	// FushallToDatabase --> For fushing all the data to Persistant Database

}
