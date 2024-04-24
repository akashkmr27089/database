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
	isPreload      bool
	trieSearch     *Trie
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

func (c *Collection) PreloadCollection(uuid string) (bool, error) {
	objectCollection, err := ReadFileAndReturnObject(uuid)
	if err != nil {
		return false, err
	}

	commands := objectCollection.Data
	if len(commands) == 0 {
		return false, nil
	}

	for _, command := range commands {
		switch command.Name {
		case InsertOneOperationName:
			c.InsertOne(*command.Key, *command.Value)
		case UpdateOneOperationName:
			c.UpdateOne(*command.Key, *command.Value)
		case DeleteOneOperationName:
			c.DeleteOne(*command.Key)
		}
	}

	c.isPreload = false
	return true, nil
}

func (c *Collection) Create(name string, updateStrategy *UpdateStrategy, preUuid *string) {
	c.trieSearch = NewTrie()
	if updateStrategy == nil {
		c.updateStrategy = BatchUpdateStrategy
	}
	c.name = name
	c.Documents = make(map[string]string, 0)
	if preUuid != nil {
		c.uuid = *preUuid

		// Run the Preload logic -->
		c.isPreload = true
		_, err := c.PreloadCollection(*preUuid)
		if err != nil {
			return
		}
	} else {
		c.uuid = uuid.New().String()

		// Add Element in the Operation Array
		operation := Operation{
			Name:  CreateCollectionOperationName,
			Value: &name,
		}
		c.incrementFlashCounter(operation)
	}
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
		isInserted = true
		c.trieSearch.Insert(key, value)

		// Add Element in the Operation Array
		if !c.isPreload {
			c.clearkeyList()
			operation := Operation{
				Name:  InsertOneOperationName,
				Key:   &key,
				Value: &value,
			}
			c.incrementFlashCounter(operation)
		}
	}

	return isInserted
}

func (c *Collection) UpdateOne(key string, value string) bool {
	_, ok := c.Documents[key]
	var isUpdated bool = false
	if ok {
		c.Documents[key] = value
		isUpdated = true
		c.trieSearch.Update(key, value)

		// Add Element in the Operation Array
		if !c.isPreload {
			c.clearkeyList()
			operation := Operation{
				Name:  UpdateOneOperationName,
				Key:   &key,
				Value: &value,
			}
			c.incrementFlashCounter(operation)
		}
	}

	return isUpdated
}

func (c *Collection) DeleteOne(key string) bool {
	_, ok := c.Documents[key]
	var isDeleted bool = false
	if !ok {
		delete(c.Documents, key)
		isDeleted = true
		c.trieSearch.Delete(key)

		// Add Element in the Operation Array
		if !c.isPreload {
			c.clearkeyList()
			operation := Operation{
				Name: DeleteOneOperationName,
				Key:  &key,
			}
			c.incrementFlashCounter(operation)
		}
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
	uuid := "c18ddbcf-6d63-483f-bc53-25dcaedf89bd"
	c.Create("Test", nil, &uuid)

	c.InsertOne("test332", "val")
	c.InsertOne("test1", "val")
	c.InsertOne("test12", "val2")
	c.InsertOne("test3", "val34")
	c.UpdateOne("test3", "val232233")
	c.InsertOne("test13", "val")
	c.InsertOne("test123", "val2")
	c.InsertOne("test33", "val34")

	values := c.GetKeysByParitialByRegex("*test1*")
	for _, value := range values {
		fmt.Println(value)
	}

	g := c.trieSearch.SearchPartial("test3")
	fmt.Print(g)
	c.Close()

	// FushallToDatabase --> For fushing all the data to Persistant Database

}
