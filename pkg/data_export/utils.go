package internal_data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ClubCedille/hackqc2024/pkg/help"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
)

type GeoJSONFeatureCollection struct {
	Type     string          `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   mapobject.Geometry     `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

func ConvertHelpsToGeoJSON(helps []*help.Help) ([]byte, error) {
	features := make([]GeoJSONFeature, 0, len(helps))

	for _, help := range helps {
		properties := make(map[string]interface{})
		properties["_id"] = help.Id
		properties["need_help"] = help.NeedHelp
		properties["how_to_help"] = help.HowToHelp
		properties["how_to_use_help"] = help.HowToUseHelp
		properties["name"] = help.MapObject.Name
		properties["description"] = help.MapObject.Description
		properties["category"] = help.MapObject.Category
		properties["tags"] = help.MapObject.Tags
		properties["date"] = help.MapObject.Date

		feature := GeoJSONFeature{
			Type:       "Feature",
			Geometry:   help.MapObject.Geometry,
			Properties: properties,
		}

		features = append(features, feature)
	}

	geoJSON := GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	return json.Marshal(geoJSON)
}

func fileSizeInMB(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s", err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %s", err)
	}

	size := fi.Size()
	
	if size < 1024 * 1024 {
		return "1" // lower than one Mo. API returns 409 error for values lower than 1.
	}
	return fmt.Sprintf("%.2f", float64(size)/1024/1024)
}