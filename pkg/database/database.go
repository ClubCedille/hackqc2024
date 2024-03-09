package database

import (
	"fmt"

	"github.com/ostafen/clover/v2"
)

const HackQcCollection = "HackQcCollection"

const AccountCollection = "AccountCollection"
const EventCollection = "EventCollection"
const HelpCollection = "HelpCollection"

func createCollectionIfNotExists(collectionName string, db *clover.DB) error {
	exists, err := db.HasCollection(collectionName)
	if err != nil {
		return err
	}

	if !exists {
		err = db.CreateCollection(collectionName)
		if err != nil {
			return err
		}
	}

	return nil
}

func InitDatabase() (*clover.DB, error) {
	db, err := clover.Open("clover-db")
	if err != nil {
		fmt.Print("Failed to open database")
		return nil, err
	}

	createCollectionIfNotExists(AccountCollection, db)
	//We reload the events at startup for now
	db.DropCollection(EventCollection)
	db.CreateCollection(EventCollection)
	createCollectionIfNotExists(HelpCollection, db)

	return db, nil
}
