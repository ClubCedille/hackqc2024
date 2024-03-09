package help

import (
	"github.com/ClubCedille/hackqc2024/pkg/database"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	uuid "github.com/satori/go.uuid"
)

type Help struct {
	Id           string              `clover:"_id"`
	MapObject    mapobject.MapObject `clover:"map_object"`
	ContactInfos string              `clover:"contact_infos"`
	NeedHelp     bool                `clover:"need_help"`
	HowToHelp    string              `clover:"how_to_help"`
	HowToUseHelp string              `clover:"how_to_use_help"`
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
	existingHelp, err := GetHelpById(db, help.Id)
	if err != nil {
		return err
	}

	existingHelp.MapObject = help.MapObject
	existingHelp.ContactInfos = help.ContactInfos
	existingHelp.NeedHelp = help.NeedHelp
	existingHelp.HowToHelp = help.HowToHelp
	existingHelp.HowToUseHelp = help.HowToUseHelp

	return UpdateHelp(db, existingHelp)
}

func DeleteHelpById(db *clover.DB, helpId string) (bool, error) {
	err := db.DeleteById(database.HelpCollection, helpId)
	if err != nil {
		return false, err
	}
	return true, nil
}
