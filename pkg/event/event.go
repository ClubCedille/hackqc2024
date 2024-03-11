package event

import (
	"fmt"
	"slices"

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
	ExternalId  string              `json:"external_id" clover:"external_id"`
	DangerLevel DangerLevel         `json:"danger_level" clover:"danger_level"`
	UrgencyType UrgencyType         `json:"urgency_type" clover:"urgency_type"`
	MapObject   mapobject.MapObject `json:"map_object" clover:"map_object"`
}

func (event *Event) GetUrgencyTypeString() string {
	switch event.UrgencyType {
	case Futur:
		return "Futur"
	case Present:
		return "Présent"
	case Past:
		return "Passé"
	default:
		return ""
	}
}

func (event *Event) GetUrgencyColor() string {
	switch event.UrgencyType {
	case Futur:
		return "orange"
	case Present:
		return "red"
	case Past:
		return "green"
	default:
		return ""
	}
}

func (event *Event) GetDangerLevelString() string {
	switch event.DangerLevel {
	case High:
		return "Élevé"
	case Medium:
		return "Modéré"
	case Low:
		return "Faible"
	default:
		return ""
	}
}

func (event *Event) GetDangerColor() string {
	switch event.DangerLevel {
	case High:
		return "red"
	case Medium:
		return "orange"
	case Low:
		return "green"
	default:
		return ""
	}
}

func (event *Event) GetCategoryColor() string {
	switch event.MapObject.Category {
	case "Pluie":
		return "blue"
	case "Neige":
		return "light-gray"
	case "Vent":
		return "gray"
	case "Onde de tempête":
		return "purple"
	default:
		return "light-gray"
	}
}

func (event *Event) GetCategoryEmoji() string {
	switch event.MapObject.Category {
	case "Pluie":
		return "🌧️"
	case "Neige":
		return "❄️"
	case "Vent":
		return "💨"
	case "Onde de tempête":
		return "🌊"
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

func GetEventWithFilters(conn *clover.DB, filters map[string][]string, requireGeoJson bool) ([]*Event, error) {
	filterQuery := query.NewQuery(database.EventCollection)

	for k, v := range filters {
		filterQuery = filterQuery.MatchFunc(func(doc *document.Document) bool {
			return slices.Contains(v, fmt.Sprint(doc.Get(k))) && (!requireGeoJson || doc.Get("map_object.geometry.coordinates") != nil)
		})
	}

	docs, err := conn.FindAll(filterQuery)
	if err != nil {
		return nil, err
	}

	events, _ := GetEventFromDocuments(docs)

	return events, nil
}

func EventExistsByExternalId(conn *clover.DB, externalId string) (bool, error) {
	return conn.Exists(query.NewQuery(database.EventCollection).Where(query.Field("external_id").Eq(externalId)))
}

func GetEventByExternalId(conn *clover.DB, externalId string) (Event, error) {
	docs, err := conn.FindFirst(query.NewQuery(database.EventCollection).Where(query.Field("external_id").Eq(externalId)))
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

	events, _ := GetEventFromDocuments(docs)

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

func GetEventFromDocuments(docs []*document.Document) ([]*Event, error) {
	var events []*Event
	for _, d := range docs {
		var event Event
		d.Unmarshal(&event)
		events = append(events, &event)
	}

	return events, nil
}
