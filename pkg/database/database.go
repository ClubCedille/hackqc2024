package database

import (
	"fmt"

	"github.com/ostafen/clover/v2"
)

const HackQcCollection = "HackQcCollection"

const AccountCollection = "AccountCollection"
const EventCollection = "EventCollection"
const HelpCollection = "HelpCollection"
const WatermarkCollection = "WatermarkCollection"
const CommentCollection = "CommentCollection"

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
	createCollectionIfNotExists(EventCollection, db)
	createCollectionIfNotExists(HelpCollection, db)
	createCollectionIfNotExists(WatermarkCollection, db)
	createCollectionIfNotExists(CommentCollection, db)

	return db, nil
}

func ExportDatabase(db *clover.DB) error {
	db.ExportCollection(AccountCollection, "account.json")
	db.ExportCollection(EventCollection, "event.json")
	db.ExportCollection(HelpCollection, "help.json")
	db.ExportCollection(WatermarkCollection, "watermark.json")
	db.ExportCollection(CommentCollection, "comment.json")

	return nil
}

func ImportDatabase(db *clover.DB) error {
	err := db.ImportCollection("account.json", AccountCollection)
	if err != nil {
		return err
	}

	err = db.ImportCollection("event.json", EventCollection)
	if err != nil {
		return err
	}

	err = db.ImportCollection("help.json", HelpCollection)
	if err != nil {
		return err
	}

	err = db.ImportCollection("watermark.json", WatermarkCollection)
	if err != nil {
		return err
	}

	err = db.ImportCollection("comment.json", CommentCollection)
	if err != nil {
		return err
	}

	return nil
}
