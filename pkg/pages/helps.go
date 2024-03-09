package pages

import (
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
)

func HelpPage(c *gin.Context, db *clover.DB) {
	docs, err := db.FindAll(query.NewQuery(database.HackQcCollection))
	if err != nil {
		log.Println("Error fetching help cards:", err)
		return
	}

	c.HTML(http.StatusOK, "cards/helpCard.html", gin.H{
		"HelpCards": docs,
	})
}

func CreateHelp(c *gin.Context, db *clover.DB) {
	var data help.Help
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := help.CreateHelp(db, data)
	if err != nil {
		log.Println("Error creating help:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Help created successfully")
	c.Redirect(http.StatusSeeOther, "/helps")
}

func UpdateHelp(c *gin.Context, db *clover.DB) {
	var data help.Help
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := help.UpdateHelp(db, data)
	if err != nil {
		log.Println("Error updating help:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Help updated successfully")
	c.Redirect(http.StatusSeeOther, "/helps")
}

func DeleteHelp(c *gin.Context, db *clover.DB) {
	var data help.Help
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	err := help.DeleteHelp(db, data)
	if err != nil {
		log.Println("Error deleting help:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Help deleted successfully")
	c.Redirect(http.StatusSeeOther, "/helps")
}
