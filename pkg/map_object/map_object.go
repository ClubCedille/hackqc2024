package mapobject

import (
	"time"
)

type MapObject struct {
	Geometry    Geometry  `json:"geometry" clover:"geometry"`
	Name        string    `json:"name" clover:"name"`
	Description string    `json:"description" clover:"description"`
	Category    string    `json:"category" clover:"category"`
	Tags        []string  `json:"tags" clover:"tags"`
	Date        time.Time `json:"date" clover:"date"`
	AccountId   string    `json:"account_id" clover:"account_id"`
}

type Geometry struct {
	GeomType    string    `json:"type" clover:"type"`
	Coordinates []float64 `json:"coordinates" clover:"coordinates"`
}
