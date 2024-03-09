package mapobject

import (
	"time"
)

type MapObject struct {
	Coordinates string    `json:"coordinates" clover:"coordinates"` //GeoJSON: Point
	Polygon     string    `json:"polygon" clover:"polygon"`         //GeoJSON: Polygon
	Name        string    `json:"name" clover:"name"`
	Description string    `json:"description" clover:"description"`
	Category    string    `json:"category" clover:"category"`
	Tags        []string  `json:"tags" clover:"tags"`
	Date        time.Time `json:"date" clover:"date"`
	AccountId   string    `json:"account_id" clover:"account_id"`
}
