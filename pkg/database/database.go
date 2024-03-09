package database

import (
	"fmt"

	"github.com/ostafen/clover/v2"
)

const HackQcCollection = "HackQcCollection"

const AccountCollection = "AccountCollection"
const MapObjectCollection = "MapObjectCollection"
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
	createCollectionIfNotExists(MapObjectCollection, db)
	createCollectionIfNotExists(EventCollection, db)
	createCollectionIfNotExists(HelpCollection, db)

	return db, nil
}
