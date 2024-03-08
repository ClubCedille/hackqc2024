package event

import (
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
)

type UrgencyType int
type DangerLevel int

const (
	Futur UrgencyType = iota
	Present
	Past
)

const (
	High DangerLevel = iota
	Medium
	Low
)

type Event struct {
	DangerLevel DangerLevel `clover:"danger_level"`
	UrgencyType UrgencyType `clover:"urgency_type"`
	MapObjectId string      `clover:"map_object_id"`
}

func CreateEvent(conn *clover.DB, event Event) error {
	eventDoc := document.NewDocumentOf(event)
	err := conn.Insert(database.HackQcCollection, eventDoc)
	if err != nil {
		return err
	}

	return nil
}
