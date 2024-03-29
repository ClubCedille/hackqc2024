package help

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	uuid "github.com/satori/go.uuid"
)

type Help struct {
	Id                   string              `json:"_id" clover:"_id"`
	MapObject            mapobject.MapObject `json:"map_object" clover:"map_object"`
	ContactInfos         string              `json:"contact_infos" clover:"contact_infos"`
	NeedHelp             bool                `json:"need_help" clover:"need_help"`
	HowToHelp            string              `json:"how_to_help" clover:"how_to_help"`
	HowToUseHelp         string              `json:"how_to_use_help" clover:"how_to_use_help"`
	EventId              string              `json:"event_id" clover:"event_id"`
	Exported             bool                `json:"exported" clover:"exported"`
	Modified             bool                `json:"modified" clover:"modified"`
	DerniereModification time.Time           `json:"derniere_modification" clover:"derniere_modification"`
}

func GetHelpById(db *clover.DB, helpId string) (Help, error) {
	docs, err := db.FindFirst(query.NewQuery(database.HelpCollection).Where(query.Field("_id").Eq(helpId)))
	if err != nil {
		return Help{}, err
	}

	help := Help{}
	docs.Unmarshal(&help)

	return help, nil
}

func GetAllHelps(db *clover.DB) ([]*Help, error) {
	docs, err := db.FindAll(query.NewQuery(database.HelpCollection))
	if err != nil {
		return nil, err
	}

	var helps []*Help
	for _, d := range docs {
		var help Help
		d.Unmarshal(&help)
		helps = append(helps, &help)
	}

	return helps, nil
}

func GetAllHelpsByAccountId(db *clover.DB, accountId string) ([]*Help, error) {
	docs, err := db.FindAll(query.NewQuery(database.HelpCollection).Where(query.Field("map_object.account_id").Eq(accountId)))
	if err != nil {
		return nil, err
	}

	var helps []*Help
	for _, d := range docs {
		var help Help
		d.Unmarshal(&help)
		helps = append(helps, &help)
	}

	return helps, nil
}

func GetHelpWithFilters(conn *clover.DB, filters map[string][]string, requireGeoJson bool) ([]*Help, error) {
	filterQuery := query.NewQuery(database.HelpCollection)
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

	helps, _ := GetHelpFromDocuments(docs)

	return helps, nil
}

func SearchHelps(conn *clover.DB, filters map[string][]string, requireGeoJson bool) ([]*Help, error) {
	filterQuery := query.NewQuery(database.HelpCollection)
	filterQuery = filterQuery.MatchFunc(func(doc *document.Document) bool {
		geoRes := !requireGeoJson || doc.Get("map_object.geometry.coordinates") != nil

		filterRes := true
		for k, v := range filters {
			if k == "map_object.category" {
				continue
			}
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

	events, _ := GetHelpFromDocuments(docs)

	return events, nil
}

func CreateHelp(db *clover.DB, help Help) error {
	help.Id = uuid.NewV4().String()
	helpDoc := document.NewDocumentOf(help)
	err := db.Insert(database.HelpCollection, helpDoc)
	if err != nil {
		return err
	}

	return nil
}

func UpdateHelp(db *clover.DB, help Help) error {
	err := db.UpdateById(database.HelpCollection, help.Id, func(doc *document.Document) *document.Document {
		doc.Set("map_object", help.MapObject)
		doc.Set("contact_infos", help.ContactInfos)
		doc.Set("need_help", help.NeedHelp)
		doc.Set("how_to_help", help.HowToHelp)
		doc.Set("how_to_use_help", help.HowToUseHelp)
		doc.Set("exported", help.Exported)
		doc.Set("derniere_modification", help.MapObject.Date)
		doc.Set("modified", true)
		return doc
	})

	if err != nil {
		return err
	}
	return nil
}

func DeleteHelpById(db *clover.DB, helpId string) error {
	err := db.DeleteById(database.HelpCollection, helpId)
	if err != nil {
		return err
	}
	return nil
}

func GetHelpFromDocuments(docs []*document.Document) ([]*Help, error) {
	var helps []*Help
	for _, d := range docs {
		var help Help
		d.Unmarshal(&help)
		helps = append(helps, &help)
	}

	return helps, nil
}

func (m *Help) GetModificationDateString() string {
	if !m.Modified {
		return ""
	}

	loc, err := time.LoadLocation("America/Montreal")
	if err != nil {
		return m.DerniereModification.Format("2 Jan 2006 à 15:04")
	}

	var month string
	switch m.DerniereModification.In(loc).Format("Jan") {
	case "Jan":
		month = "janvier"
	case "Feb":
		month = "février"
	case "Mar":
		month = "mars"
	case "Apr":
		month = "avril"
	case "May":
		month = "mai"
	case "Jun":
		month = "juin"
	case "Jul":
		month = "juillet"
	case "Aug":
		month = "août"
	case "Sep":
		month = "septembre"
	case "Oct":
		month = "octobre"
	case "Nov":
		month = "novembre"
	case "Dec":
		month = "décembre"
	}
	return m.DerniereModification.Format("2") + " " + month + " " + m.DerniereModification.Format("2006 à 15:04")
}

func (help *Help) FlipCoords() {
	if help.MapObject.Geometry.GeomType == "Point" {
		tmp := help.MapObject.Geometry.Coordinates[0]
		help.MapObject.Geometry.Coordinates[0] = help.MapObject.Geometry.Coordinates[1]
		help.MapObject.Geometry.Coordinates[1] = tmp
	}
}
