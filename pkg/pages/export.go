package pages

import (
	"log"
	"net/http"
	"os"

	data_export "github.com/ClubCedille/hackqc2024/pkg/data_export"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func SubmitHelpsToDC(c *gin.Context, db *clover.DB) {

	filePath := "tmp/soumission-aide.json"
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

	err := data_export.PostJsonHelpsToDQ(apiKey, jeuDeDonnees, filePath)
	if err != nil {
		log.Printf("Error posting help events to Données Québec: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post help data"})
		return
	}

	// helps, err := help.GetAllHelps(db)
	// if err != nil {
	// 	log.Println("Error fetching events:", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch help data"})
	// 	return
	// }

	// err = data_export.PostGeoJsonHelpsToDQ(apiKey, "1eba7e31-a048-47fa-ab28-d2aa0cdec51d", helps)
	// if err != nil {
	// 	log.Printf("Error posting help events to Données Québec: %s", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to post help data"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{"status": "submitted"})
}