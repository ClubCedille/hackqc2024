package mapobject

import (
	"time"
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
