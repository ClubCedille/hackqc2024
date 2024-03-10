package pages

import (
	"log"
	"net/http"
	"os"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/internal_data"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func SubmitEvents(c *gin.Context, db *clover.DB) {

	filePath := "tmp/events.json"
	db.ExportCollection("EventCollection", filePath)

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("API_KEY environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "API key not set"})
		return
	}
	
	err := internal_data.PostJsonEventsToDQ(apiKey, "1eba7e31-a048-47fa-ab28-d2aa0cdec51d", filePath)
	if err != nil {
		log.Printf("Error posting events to Données Québec: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post event data"})
		return
	}

	events, err := event.GetAllEvents(db)
	if err != nil {
		log.Println("Error fetching events:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch event data"})
		return
	}

	err = internal_data.PostGeoJsonEventsToDQ(apiKey, "1eba7e31-a048-47fa-ab28-d2aa0cdec51d", events)
	if err != nil {
		log.Printf("Error posting events to Données Québec: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post event data"})
		return
	}	

	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}