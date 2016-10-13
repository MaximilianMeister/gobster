// Package feed provides a key value storage for the data
package feed

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/xyproto/simplebolt"
)

// Open establishes and returns a database connection
func Open() (*simplebolt.Database, error) {
	db, err := simplebolt.New("gobster.db")

	if err != nil {
		return db, err
	}

	return db, err
}

// Get retrieves a random value (list entry) from key/value storage
func Get(listName string) (string, error) {
	allEntries, err := GetAll(listName)
	if err != nil {
		fmt.Printf("Could not get list! %s", err)
		return "", err
	}

	// seed to enforce real randomness
	rand.Seed(time.Now().UnixNano())

	randomEntry := ""
	if len(allEntries) > 0 {
		randomEntry = allEntries[rand.Intn(len(allEntries))]
	}

	return randomEntry, nil
}

func GetAll(listName string) ([]string, error) {
	db, err := Open()
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	list, err := simplebolt.NewList(db, listName)
	if err != nil {
		fmt.Printf("Could not create a list! %s", err)
		return []string{}, err
	}

	allEntries, err := list.GetAll()
	if err != nil {
		fmt.Printf("Could not get an item from the list! %s", err)
		return []string{}, err
	}

	// remove the list again if it is empty
	if len(allEntries) == 0 {
		err := list.Remove()
		if err != nil {
			fmt.Printf("Could not remove the list! %s", err)
		}
	}

	return allEntries, nil
}

// Set writes a value (list entry) to a key/value storage
func Set(listName string, value string) error {
	db, err := Open()
	if err != nil {
		return err
	}
	defer db.Close()

	list, err := simplebolt.NewList(db, listName)
	if err != nil {
		fmt.Printf("Could not create a list! %s", err)
		return err
	}

	if err := list.Add(value); err != nil {
		fmt.Printf("Could not add an item to the list! %s", err)
		return err
	}

	return nil
}
