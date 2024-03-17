package event

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ClubCedille/hackqc2024/pkg/notifications"
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
	Id                  string              `json:"_id" clover:"_id"`
	ExternalId          string              `json:"external_id" clover:"external_id"`
	DangerLevel         DangerLevel         `json:"danger_level" clover:"danger_level"`
	UrgencyType         UrgencyType         `json:"urgency_type" clover:"urgency_type"`
	Subscribers         []string            `json:"subscribers" clover:"subscribers"`
	MapObject           mapobject.MapObject `json:"map_object" clover:"map_object"`
	MunicipalityPolygon [][][]float64       `json:"municipality_polygon" clover:"municipality_polygon"`
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
	filterQuery = filterQuery.MatchFunc(func(doc *document.Document) bool {
		result := !requireGeoJson || doc.Get("map_object.geometry.coordinates") != nil
		for k, v := range filters {
			result = result && (!doc.Has(k) || slices.Contains(v, fmt.Sprint(doc.Get(k))))
		}
		return result
	})

	docs, err := conn.FindAll(filterQuery)
	if err != nil {
		return nil, err
	}

	events, _ := GetEventFromDocuments(docs)

	return events, nil
}

func SearchEvents(conn *clover.DB, filters map[string][]string, requireGeoJson bool) ([]*Event, error) {
	filterQuery := query.NewQuery(database.EventCollection)
	filterQuery = filterQuery.MatchFunc(func(doc *document.Document) bool {
		geoRes := !requireGeoJson || doc.Get("map_object.geometry.coordinates") != nil

		filterRes := true
		for k, v := range filters {
			docKeyValue := doc.Get(k)
			hasKey := doc.Has(k)
			hasValue := slices.Contains(v, fmt.Sprint(docKeyValue))
			filterRes = filterRes && (!hasKey || hasValue)
		}

		// If there are search terms, check if the document contains any of them
		searchRes := false
		searchTerms := filters["search"]
		for _, search := range searchTerms {
			if search == "" {
				searchRes = true
				continue
			}
			for k := range doc.AsMap() {
				searchRes = searchRes || (strings.Contains(strings.ToLower(fmt.Sprint(doc.Get(k))), strings.ToLower(search)))
			}
		}

		finalRes := ((geoRes && filterRes) && (searchRes))

		return finalRes
	})

	// If we have _.sort, we need to sort the results by that field
	if filters["_.sort"] != nil {
		var sortOrder int
		if filters["_.sortOrder"][0] == "1" {
			sortOrder = 1 // Ascending
		} else {
			sortOrder = -1 // Descending
		}
		sortField := filters["_.sort"][0]
		filterQuery = filterQuery.Sort(query.SortOption{Field: sortField, Direction: sortOrder})
	} else {
		filterQuery = filterQuery.Sort(query.SortOption{Field: "map_object.date", Direction: -1})
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

func GetAllEventsByAccountId(conn *clover.DB, accountId string) ([]*Event, error) {
	docs, err := conn.FindAll(query.NewQuery(database.EventCollection).Where(query.Field("map_object.account_id").Eq(accountId)))
	if err != nil {
		return nil, err
	}

	events, err := GetEventFromDocuments(docs)
	if err != nil {
		return nil, err
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

	message := fmt.Sprintf("Alert: Un nouvel évènement de type %s a été signalé près de vous", event.MapObject.Category)
	err = notifications.NotifyNearby(conn, message, event.MapObject.Geometry)

	return nil
}

func UpdateEvent(conn *clover.DB, event Event) error {
	err := conn.UpdateById(database.EventCollection, event.Id, func(doc *document.Document) *document.Document {
		doc.Set("danger_level", event.DangerLevel)
		doc.Set("urgency_type", event.UrgencyType)
		doc.Set("subscribers", event.Subscribers)
		doc.Set("map_object", event.MapObject)
		return doc
	})
	if err != nil {
		return err
	}

	return nil
}

func DeleteEventById(conn *clover.DB, eventId string) error {
	err := conn.DeleteById(database.EventCollection, eventId)
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

func (event *Event) FlipCoords() {
	if event.MapObject.Geometry.GeomType == "Point" {
		tmp := event.MapObject.Geometry.Coordinates[0]
		event.MapObject.Geometry.Coordinates[0] = event.MapObject.Geometry.Coordinates[1]
		event.MapObject.Geometry.Coordinates[1] = tmp
	}
}
