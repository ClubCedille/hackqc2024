package pages

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	data_export "github.com/ClubCedille/hackqc2024/pkg/data_export"
	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func SubmitHelpsToDC(c *gin.Context, db *clover.DB) {

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

	updatedHelp, err := prepareHelpDataForExport(filePath, events)
	if err != nil {
		log.Printf("Error updating external source linked to help: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update help data"})
		return
	}

	err = data_export.PostJsonHelpsToDQ(apiKey, jeuDeDonnees, filePath)
	if err != nil {
		log.Printf("Error posting help events to Données Québec: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post help data"})
		return
	}

	if err != nil {
		log.Println("Error converting help doc to helps:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update help data"})
		return
	}

	err = data_export.PostGeoJsonHelpsToDQ(apiKey, jeuDeDonnees, updatedHelp)
	if err != nil {
		log.Printf("Error posting help events to Données Québec: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post help data"})
		return
	}

	log.Println("Helps submitted to Données Québec")
}

func prepareHelpDataForExport(filePath string, linkedEvents []*event.Event) ([]map[string]interface{}, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        log.Printf("Failed to read the file: %v", err)
        return nil, err
    }

    var docs []map[string]interface{}
    if err := json.Unmarshal(data, &docs); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON data: %w", err)
    }

    for _, doc := range docs {
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
        }
    }

    updatedData, err := json.Marshal(docs)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal the updated JSON data: %w", err)
    }

    if err := os.WriteFile(filePath, updatedData, 0644); err != nil {
        return nil, fmt.Errorf("failed to write the updated JSON back to the file: %w", err)
    }

    return docs, nil
}
