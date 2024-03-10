package pages

import (
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
)

func GridPage(c *gin.Context, db *clover.DB) {

	// Get all events and helps
	var events []*event.Event
	events, _ = event.GetAllEvents(db)

	var helps []*help.Help
	helps, _ = help.GetAllHelps(db)

	c.HTML(http.StatusOK, "grid/index.html", gin.H{
		"Events": events,
		"Helps":  helps,
	})
}

func GridSearch(c *gin.Context, db *clover.DB) {

	// Search query
	searchTerm := c.Query("search")

	if searchTerm == "" {
		// No search query, show all events and helps
		var events []*event.Event
		events, _ = event.GetAllEvents(db)

		var helps []*help.Help
		helps, _ = help.GetAllHelps(db)

		c.HTML(http.StatusOK, "components/grid", gin.H{
			"Events": events,
			"Helps":  helps,
		})
		return
	}

	// Search for events
	docs, err := db.FindAll(query.NewQuery(database.EventCollection).Where(query.Field("map_object.name").Like(searchTerm).Or(query.Field("map_object.description").Like(searchTerm).Or(query.Field("map_object.category").Like(searchTerm)).Or(query.Field("map_object.tags").Contains(searchTerm)))))
	if err != nil {
		log.Println("Error fetching events:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	events, _ := event.GetEventFromDocuments(docs)

	// Search for helps
	docs, err = db.FindAll(query.NewQuery(database.HelpCollection).Where(query.Field("map_object.name").Like(searchTerm).Or(query.Field("map_object.description").Like(searchTerm).Or(query.Field("map_object.category").Like(searchTerm)).Or(query.Field("map_object.tags").Contains(searchTerm)))))
	if err != nil {
		log.Println("Error fetching helps:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	helps, _ := help.GetHelpFromDocuments(docs)

	c.HTML(http.StatusOK, "components/grid", gin.H{
		"Events": events,
		"Helps":  helps,
	})
	return
}
