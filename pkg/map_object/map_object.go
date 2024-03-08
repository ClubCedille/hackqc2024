package mapobject

import (
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
)

type MapObject struct {
	Coordinates string    `clover:"coordinates"` //GeoJSON: Point
	Polygon     string    `clover:"polygon"`     //GeoJSON: Polygon
	Name        string    `clover:"name"`
	Description string    `clover:"description"`
	Category    string    `clover:"category"`
	Tags        []string  `clover:"tags"`
	Date        time.Time `clover:"date"`
	AccountId   string    `clover:"account_id"`
}

func CreateMapObject(conn *clover.DB, mapObject MapObject) error {
	mapObjectDoc := document.NewDocumentOf(mapObject)
	err := conn.Insert(database.HackQcCollection, mapObjectDoc)
	if err != nil {
		return err
	}

	return nil
}
