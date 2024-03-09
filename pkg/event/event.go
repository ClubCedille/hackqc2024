package event

import (
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
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

func (event *Event) GetUrgencyTypeString() string {
	switch event.UrgencyType {
	case Futur:
		return "Futur"
	case Present:
		return "Present"
	case Past:
		return "Past"
	default:
		return ""
	}
}

func (event *Event) GetDangerLevelString() string {
	switch event.DangerLevel {
	case High:
		return "High"
	case Medium:
		return "Medium"
	case Low:
		return "Low"
	default:
		return ""
	}
}

func GetEventById(conn *clover.DB, eventId string) (Event, error) {
	docs, err := conn.FindAll(query.NewQuery(database.EventCollection).Where(query.Field("_id").Eq(eventId)))
	if err != nil {
		return Event{}, err
	}

	return Event{
		DangerLevel: DangerLevel(docs[0].Get("danger_level").(int)),
		UrgencyType: UrgencyType(docs[0].Get("urgency_type").(int)),
		MapObjectId: docs[0].Get("map_object_id").(string),
	}, nil
}

func GetAllEvents(conn *clover.DB) ([]*Event, error) {
	docs, err := conn.FindAll(query.NewQuery(database.EventCollection))
	if err != nil {
		return nil, err
	}

	var events []*Event
	for _, d := range docs {
		var event Event
		d.Unmarshal(&event)
		events = append(events, &event)
	}

	return events, nil
}

func CreateEvent(conn *clover.DB, event Event) error {
	eventDoc := document.NewDocumentOf(event)
	err := conn.Insert(database.EventCollection, eventDoc)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEvent(conn *clover.DB, event Event) error {
	return nil
}

func DeleteEvent(conn *clover.DB, event Event) error {
	return nil
}
