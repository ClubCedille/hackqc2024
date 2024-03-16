package pages

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	data_export "github.com/ClubCedille/hackqc2024/pkg/data_export"
	"github.com/ClubCedille/hackqc2024/pkg/event"
	circletopolygon "github.com/chrusty/go-circle-to-polygon"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func SubmitHelpsToDC(c *gin.Context, db *clover.DB, helpIds []string) {

	filePath := "tmp/soumissions-aide.json"
	db.ExportCollection("HelpCollection", filePath)

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("API_KEY environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API key not set"})
		return
	}


	jeuDeDonnees := os.Getenv("JEU_DE_DONNEES")
	if jeuDeDonnees == "" {
		log.Println("JEU_DE_DONNEES environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Jeu de données not set"})
		return
	}

	events, err := event.GetAllEvents(db)
	if err != nil {
		log.Println("Error fetching events:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch help data"})
		return
	}

	updatedHelp, err := prepareHelpDataForExport(filePath, events, helpIds)
	if err != nil {
		log.Printf("Error updating external source linked to help: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update help data"})
		return
	}

	err = data_export.PostJsonHelpsToDQ(apiKey, jeuDeDonnees, filePath)
	if err != nil {
		log.Printf("Error posting json help events to Données Québec: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post help data"})
		return
	}

	err = data_export.PostGeoJsonHelpsToDQ(apiKey, jeuDeDonnees, updatedHelp)
	if err != nil {
		log.Printf("Error posting geojson help events to Données Québec: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post help data"})
		return
	}

	log.Println("Helps submitted to Données Québec")
}

func prepareHelpDataForExport(filePath string, linkedEvents []*event.Event, helpIds []string) ([]map[string]interface{}, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        log.Printf("Failed to read the file: %v", err)
        return nil, err
    }

    var docs []map[string]interface{}
    if err := json.Unmarshal(data, &docs); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
    }

	filteredDocs := []map[string]interface{}{}
    for _, doc := range docs {
        if id, ok := doc["Id"].(string); ok {
            if contains(helpIds, id) {
                filteredDocs = append(filteredDocs, doc)
            }
        }
    }

    for _, doc := range filteredDocs {
        if eventID, ok := doc["event_id"].(string); ok {
            for _, event := range linkedEvents {
                if event.Id == eventID {
                    doc["source_externe_linked"] = event.ExternalId
                    doc["categorie_catastrophe"] = event.MapObject.Category
                    break
                }
            }
            delete(doc, "event_id")
            delete(doc, "contact_infos")
        }
        if mapObject, ok := doc["map_object"].(map[string]interface{}); ok {
            delete(mapObject, "account_id")
            delete(mapObject, "Id")

			if geometry, ok := mapObject["geometry"].(map[string]interface{}); ok {
				if coordinates, ok := geometry["coordinates"].([]interface{}); ok && len(coordinates) == 2 {
					latitude, latOk := coordinates[1].(float64)
					longitude, longOk := coordinates[0].(float64)
		
					if latOk && longOk {
						radius := 1000
						edges := 10
		
						geoJSON := convertToCirclePolygon(&latitude, &longitude, &radius, &edges)

						var geoJSONObject map[string]interface{}
						if err := json.Unmarshal([]byte(geoJSON), &geoJSONObject); err != nil {
							log.Printf("Error parsing the GeoJSON string: %v", err)
							return nil, err
						} else {
							mapObject["geometry"] = geoJSONObject
						}
					}
				}
			}
        }
    }
    updatedData, err := json.Marshal(filteredDocs)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal the updated JSON data: %w", err)
    }

    if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
        return nil, fmt.Errorf("failed to write the updated JSON back to the file: %w", err)
    }

    return filteredDocs, nil
}


func convertToCirclePolygon(latitude *float64, longitude *float64, radius *int, edges *int) []byte {
	// Make a circle:
	circle := &circletopolygon.Circle{
		Centre: &circletopolygon.Point{
			Latitude:  float32(*latitude),
			Longitude: float32(*longitude),
		},
		Radius: int32(*radius),
	}

	// Validate the circle:
	if err := circle.Validate(); err != nil {
		panic(err)
	}

	// Convert it to a Polygon with 10 edges:
	polygon, err := circle.ToPolygon(int(*edges))
	if err != nil {
		panic(err)
	}

	// Render as GeoJSON:
	geoJSON, err := polygon.GeoJSON()
	if err != nil {
		panic(err)
	}

	return geoJSON
}

func contains(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}