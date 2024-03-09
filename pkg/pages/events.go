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

func EventsPage(c *gin.Context, db *clover.DB) {
	events, err := event.GetAllEvents(db)
	if err != nil {
		log.Println("Error fetching event cards:", err)
		return
	}

	c.HTML(http.StatusOK, "cards/eventCard.html", gin.H{
		"EventCards": events,
	})
}

func EventHelpPage(c *gin.Context, db *clover.DB) {
	events, err := event.GetAllEvents(db)

	if err != nil {
		log.Println("Error fetching events:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	helps, err := help.GetAllHelps(db)

	if err != nil {
		log.Println("Error fetching helps:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "list/index.html", gin.H{
		"Events": events,
		"Helps":  helps,
	})
}

func SearchEventHelpPage(c *gin.Context, db *clover.DB) {
	searchTerm := c.Query("search")

	if searchTerm == "" {
		c.HTML(http.StatusOK, "list/event_list_table.html", gin.H{
			"Events": []*event.Event{},
			"Helps":  []*help.Help{},
		})
		return
	}

	docs, err := db.FindAll(query.NewQuery(database.EventCollection).Where(query.Field("map_object.name").Like(searchTerm).Or(query.Field("map_object.description").Like(searchTerm).Or(query.Field("map_object.category").Like(searchTerm)).Or(query.Field("map_object.tags").Contains(searchTerm)))))
	if err != nil {
		log.Println("Error fetching events:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	events, err := event.GetEventFromDocuments(docs)

	c.HTML(http.StatusOK, "list/event_list_table.html", gin.H{
		"Events": events,
		"Helps":  []*help.Help{},
	})
}

func CreateEvent(c *gin.Context, db *clover.DB) {
	var data event.Event
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := event.CreateEvent(db, data)
	if err != nil {
		log.Println("Error creating event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Event created successfully")
	c.Redirect(http.StatusSeeOther, "/events")
}

func UpdateEvent(c *gin.Context, db *clover.DB) {
	var data event.Event
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := event.UpdateEvent(db, data)
	if err != nil {
		log.Println("Error updating event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Event updated successfully")
	c.Redirect(http.StatusSeeOther, "/events")
}

func DeleteEvent(c *gin.Context, db *clover.DB) {
	var data event.Event
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := event.DeleteEvent(db, data)
	if err != nil {
		log.Println("Error deleting event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Event deleted successfully")
	c.Redirect(http.StatusSeeOther, "/events")
}
