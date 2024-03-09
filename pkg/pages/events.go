package pages

import (
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
)

func EventsPage(c *gin.Context, db *clover.DB) {
	docs, err := db.FindAll(query.NewQuery(database.HackQcCollection).Where(query.Field("urgency_type").Eq(1)))
	if err != nil {
		log.Println("Error fetching event cards:", err)
		return
	}

	c.HTML(http.StatusOK, "cards/eventCard.html", gin.H{
		"EventCards": docs,
	})
}

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

func EventHelpPage(c *gin.Context, db *clover.DB) {
	c.HTML(http.StatusOK, "list/index.html", nil)
}
