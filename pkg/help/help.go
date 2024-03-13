package help

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

type Help struct {
	Id           string              `json:"_id" clover:"_id"`
	MapObject    mapobject.MapObject `json:"map_object" clover:"map_object"`
	ContactInfos string              `json:"contact_infos" clover:"contact_infos"`
	NeedHelp     bool                `json:"need_help" clover:"need_help"`
	HowToHelp    string              `json:"how_to_help" clover:"how_to_help"`
	HowToUseHelp string              `json:"how_to_use_help" clover:"how_to_use_help"`
	EventId      string              `json:"event_id" clover:"event_id"`
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
