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

func (m *MapObject) GetDateString() string {
	var month string
	switch m.Date.Format("Jan") {
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
	return m.Date.Format("2") + " " + month + " " + m.Date.Format("2006 à 15:04")
}
