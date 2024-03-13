package internal_data

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type GeoJSONFeatureCollection struct {
	Type     string          `json:"type"`
	Features []GeoJSONFeature `json:"features"`
}

type GeoJSONFeature struct {
	Type       string                 `json:"type"`
	Geometry   Geometry     `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates interface{} `json:"coordinates"` 
}

func ConvertMapDocsToGeoJSON(docs []map[string]interface{}) ([]byte, error) {
    features := make([]GeoJSONFeature, 0, len(docs))

    for _, doc := range docs {
        properties := make(map[string]interface{})

        // Extract and assign values from doc. You'll need to assert the type of each value.
        if id, ok := doc["_id"].(string); ok {
            properties["_id"] = id
        }
        if needHelp, ok := doc["need_help"].(bool); ok {
            properties["need_help"] = needHelp
        }
        if howToHelp, ok := doc["how_to_help"].(string); ok {
            properties["how_to_help"] = howToHelp
        }
		if howToUseHelp, ok := doc["how_to_use_help"].(string); ok {
			properties["how_to_use_help"] = howToUseHelp
		}
		if name, ok := doc["name"].(string); ok {
			properties["name"] = name
		}
		if categorieCatastrophe, ok := doc["categorie_catastrophe"].(string); ok {
			properties["categorie_catastrophe"] = categorieCatastrophe
		}
		if sourceExterneLinked, ok := doc["source_externe_linked"].(string); ok {
			properties["source_externe_linked"] = sourceExterneLinked
		}

        if mapObject, ok := doc["map_object"].(map[string]interface{}); ok {
            if name, ok := mapObject["name"].(string); ok {
                properties["name"] = name
            }
            if description, ok := mapObject["description"].(string); ok {
                properties["description"] = description
            }
            if category, ok := mapObject["category"].(string); ok {
                properties["category"] = category
            }
            if tags, ok := mapObject["tags"].([]string); ok {
				properties["tags"] = tags
			}
			if date, ok := mapObject["date"].(string); ok {
				properties["date"] = date
			}

            if geomMap, ok := mapObject["geometry"].(map[string]interface{}); ok {
				geometry := Geometry{
                    Type:        geomMap["type"].(string),
                    Coordinates: geomMap["coordinates"],
                }
                feature := GeoJSONFeature{
                    Type:       "Feature",
                    Geometry:   geometry,
                    Properties: properties,
                }
				features = append(features, feature)
            }
        }
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