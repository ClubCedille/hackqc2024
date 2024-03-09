package pages

import (
	"encoding/json"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

type GeoJSON struct {
	Type       string              `json:"type"`
	Geometry   mapobject.Geometry  `json:"geometry"`
	Properties mapobject.MapObject `json:"properties"`
}

func MapPage(c *gin.Context, db *clover.DB) {
	mapItemsJson, err := retrieveMapItemsJson(db)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "map/index.html", gin.H{
		"MapItemsJson": mapItemsJson,
	})
}

func retrieveMapItemsJson(db *clover.DB) (string, error) {
	events, err := event.GetAllEvents(db)
	if err != nil {
		return "", err
	}

	helps, err := help.GetAllHelps(db)
	if err != nil {
		return "", err
	}

	evSize := len(events)
	helpSize := len(helps)
	mapItems := make([]GeoJSON, evSize+helpSize)

	for i, v := range events {
		mapItems[i] = GeoJSON{
			Type:       "Feature",
			Geometry:   v.MapObject.Geometry,
			Properties: v.MapObject,
		}
	}
	for i, v := range helps {
		mapItems[i+evSize] = GeoJSON{
			Type:       "Feature",
			Geometry:   v.MapObject.Geometry,
			Properties: v.MapObject,
		}
	}

	jsonValue, err := json.Marshal(mapItems)
	if err != nil {
		return "", err
	}

	return string(jsonValue), nil
}
