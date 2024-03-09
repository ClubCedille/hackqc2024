package event

import (
	"github.com/ClubCedille/hackqc2024/pkg/database"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	uuid "github.com/satori/go.uuid"
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
	Id          string              `json:"_id" clover:"_id"`
	DangerLevel DangerLevel         `json:"danger_level" clover:"danger_level"`
	UrgencyType UrgencyType         `json:"urgency_type" clover:"urgency_type"`
	MapObject   mapobject.MapObject `json:"map_object" clover:"map_object"`
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
	docs, err := conn.FindFirst(query.NewQuery(database.EventCollection).Where(query.Field("_id").Eq(eventId)))
	if err != nil {
		return Event{}, err
	}

	event := Event{}
	docs.Unmarshal(&event)

	return event, nil
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
	event.Id = uuid.NewV4().String()
	eventDoc := document.NewDocumentOf(event)
	err := conn.Insert(database.EventCollection, eventDoc)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEvent(conn *clover.DB, event Event) error {
	err := conn.UpdateById(database.EventCollection, event.Id, func(doc *document.Document) *document.Document {
		doc.Set("danger_level", event.DangerLevel)
		doc.Set("urgency_type", event.UrgencyType)
		doc.Set("map_object", event.MapObject)
		return doc
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteEvent(conn *clover.DB, event Event) error {
	err := conn.DeleteById(database.EventCollection, event.Id)
	if err != nil {
		return err
	}
	return nil
}
