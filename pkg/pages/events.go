package pages

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/event"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ClubCedille/hackqc2024/pkg/session"
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

func EventTablePage(c *gin.Context, db *clover.DB) {
	docs, err := db.FindAll(query.NewQuery(database.EventCollection).Sort(query.SortOption{Field: "map_object.date"}))

	var events []*event.Event
	events, _ = event.GetEventFromDocuments(docs)

	if err != nil {
		log.Println("Error fetching events:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "events/table.html", gin.H{
		"Events": events,
	})
}

func SearchEventTable(c *gin.Context, db *clover.DB) {
	searchTerm := c.Query("search")

	if searchTerm == "" {
		allEvents, _ := event.GetAllEvents(db)
		c.HTML(http.StatusOK, "components/event-table", gin.H{
			"Events": allEvents,
		})
		return
	}

	docs, err := db.FindAll(query.NewQuery(database.EventCollection).Where(query.Field("map_object.name").Like(searchTerm).Or(query.Field("map_object.description").Like(searchTerm).Or(query.Field("map_object.category").Like(searchTerm)).Or(query.Field("map_object.tags").Contains(searchTerm)))))
	if err != nil {
		log.Println("Error fetching events:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	events, _ := event.GetEventFromDocuments(docs)

	c.HTML(http.StatusOK, "components/event-table", gin.H{
		"Events": events,
	})
}

func CreateEvent(c *gin.Context, db *clover.DB) {
	dangerLvl, _ := strconv.Atoi(c.PostForm("danger_level"))
	urgencyType, _ := strconv.Atoi(c.PostForm("urgency_type"))

	tags := c.PostForm("map_object_tags")
	tagsArray := strings.Split(tags, ",")

	var tagsArrayString []string
	for _, v := range tagsArray {
		tag := strings.TrimSpace(v)
		tagsArrayString = append(tagsArrayString, tag)
	}

	coordinates := c.PostForm("map_object_geometry_coordinates")
	coordinatesArray := strings.Split(coordinates, ",")

	var coordinatesArrayFloat []float64
	for i := len(coordinatesArray) - 1; i >= 0; i-- {
		coords, err := strconv.ParseFloat(strings.TrimSpace(coordinatesArray[i]), 64)
		if err != nil {
			log.Println("Error parsing coordinates:", err)
			c.Status(http.StatusInternalServerError)
			return
		}
		coordinatesArrayFloat = append(coordinatesArrayFloat, coords)
	}

	data := event.Event{
		DangerLevel: event.DangerLevel(dangerLvl),
		UrgencyType: event.UrgencyType(urgencyType),
		MapObject: mapobject.MapObject{
			Name:        c.PostForm("map_object_name"),
			Description: c.PostForm("map_object_description"),
			Category:    c.PostForm("map_object_category"),
			Tags:        tagsArrayString,
			AccountId:   session.ActiveSession.AccountId,
			Date:        time.Now(),
			Geometry: mapobject.Geometry{
				GeomType:    c.PostForm("map_object_geometry_type"),
				Coordinates: coordinatesArrayFloat,
			},
		},
	}

	err := event.CreateEvent(db, data)
	if err != nil {
		log.Println("Error creating event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Event created successfully")
	c.Redirect(http.StatusSeeOther, "/map")
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
	eventID := c.Param("id")

	err := event.DeleteEventById(db, eventID)
	if err != nil {
		log.Println("Error deleting event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Event deleted successfully")
}

func GetEventDetailsAboutToBeDelete(c *gin.Context, db *clover.DB) {
	id := c.Param("id")
	event, err := event.GetEventById(db, id)
	if err != nil {
		log.Println("Error getting event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "modals/delete-event.html", gin.H{
		"Event": &event,
	})
}

func EventDetails(c *gin.Context, db *clover.DB) {
	id := c.Param("id")
	event, err := event.GetEventById(db, id)
	if err != nil {
		log.Println("Error getting event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "modals/event-details.html", gin.H{
		"Event": &event,
	})
}
