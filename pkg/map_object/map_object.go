package mapobject

import (
	"time"

	"github.com/JamesLMilner/pip-go"
	"github.com/twpayne/go-geom"
)

type MapObject struct {
	// These Properties are a hack to make the UI work. We should probably
	// have a view model or something but I am lazy and tired
	Id   string `json:"id"`
	Type string `json:"type"`

	//Normal stuff below
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

func (mapObject *MapObject) GetCategoryEmoji() string {
	switch mapObject.Category {
	case "Pluie":
		return "🌧️"
	case "Neige":
		return "❄️"
	case "Vent":
		return "💨"
	case "Onde de tempête":
		return "🌊"
	case "Hébergement":
		return "🛌"
	case "Nourriture":
		return "🍲"
	case "Transport":
		return "🚗"
	case "Coup de main":
		return "🤝"
	case "Renforcement":
		return "➕"
	default:
		return ""
	}
}

func (geometry *Geometry) AsGeomCoord() geom.Coord {
	if geometry.GeomType == "Point" {
		return geom.Coord(geometry.Coordinates[0:2])
	}
	return nil
}

func (geometry *Geometry) AsPipPolygon() *pip.Polygon {
	if geometry.GeomType == "Polygon" {
		points := []pip.Point{}
		for i := 0; i < len(geometry.Coordinates); i += 2 {
			x := geometry.Coordinates[i]
			y := geometry.Coordinates[i+1]
			points = append(points, pip.Point{X: x, Y: y})
		}
		return &pip.Polygon{Points: points}
	}
	return nil
}
