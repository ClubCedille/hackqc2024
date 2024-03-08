package database

import (
	"github.com/ostafen/clover/v2"
)

const HackQcCollection = "HackQcCollection"

func InitDatabase() (*clover.DB, error) {
	db, err := clover.Open("clover-db")
	if err != nil {
		return nil, err
	}

	hasCollection, err := db.HasCollection(HackQcCollection)
	if err != nil {
		return nil, err
	}

	if !hasCollection {
		db.CreateCollection(HackQcCollection)
	}

	return db, nil
}
