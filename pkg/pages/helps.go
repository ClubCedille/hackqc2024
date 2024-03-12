package pages

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
)

func HelpPage(c *gin.Context, db *clover.DB) {
	helps, err := help.GetAllHelps(db)
	if err != nil {
		log.Println("Error fetching help cards:", err)
		return
	}

	c.HTML(http.StatusOK, "cards/helpCard.html", gin.H{
		"HelpCards": helps,
	})
}

func CreateHelp(c *gin.Context, db *clover.DB) {
	eventName := c.PostForm("map_object_name")
	eventDescription := c.PostForm("map_object_description")
	eventCategory := c.PostForm("map_object_category")
	eventId := c.PostForm("event_id")

	// Processing tags
	tags := c.PostForm("map_object_tags")
	tagsArray := strings.Split(tags, ",")
	var tagsArrayString []string
	for _, tag := range tagsArray {
		trimmedTag := strings.TrimSpace(tag)
		if trimmedTag != "" {
			tagsArrayString = append(tagsArrayString, trimmedTag)
		}
	}

	// Processing coordinates
	coordinatesStr := c.PostForm("map_object_geometry_coordinates")
	coordinatesArray := strings.Split(coordinatesStr, ",")
	var coordinatesArrayFloat []float64
	for _, coord := range coordinatesArray {
		floatCoord, err := strconv.ParseFloat(strings.TrimSpace(coord), 64)
		if err != nil {
			log.Println("Error parsing coordinates:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coordinates format"})
			return
		}
		coordinatesArrayFloat = append(coordinatesArrayFloat, floatCoord)
	}

	contactInfos := c.PostForm("contact_infos")
	needHelp := c.PostForm("need_help") == "true"
	howToHelp := c.PostForm("how_to_help")
	howToUseHelp := c.PostForm("how_to_use_help")

	helpRequest := help.Help{
		ContactInfos: contactInfos,
		NeedHelp:     needHelp,
		HowToHelp:    howToHelp,
		HowToUseHelp: howToUseHelp,
		EventId:      eventId,
		MapObject: mapobject.MapObject{
			AccountId:   session.ActiveSession.AccountId,
			Name:        eventName,
			Description: eventDescription,
			Category:    eventCategory,
			Tags:        tagsArrayString,
			Geometry:    mapobject.Geometry{GeomType: "Point", Coordinates: coordinatesArrayFloat},
		},
	}

	err := help.CreateHelp(db, helpRequest)
	if err != nil {
		log.Println("Error submitting help request:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Help created successfully")
	c.Redirect(http.StatusSeeOther, "/map")
}

func UpdateHelp(c *gin.Context, db *clover.DB) {
	var data help.Help
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func HelpTablePage(c *gin.Context, db *clover.DB) {
	docs, err := db.FindAll(query.NewQuery(database.HelpCollection).Sort(query.SortOption{Field: "map_object.date"}))

	var helps []*help.Help
	helps, _ = help.GetHelpFromDocuments(docs)

	if err != nil {
		log.Println("Error fetching helps:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "helps/table.html", gin.H{
		"Helps": helps,
	})
}

func HelpDetails(c *gin.Context, db *clover.DB) {
	id := c.Param("id")
	help, err := help.GetHelpById(db, id)
	if err != nil {
		log.Println("Error getting help:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "modals/help-details.html", gin.H{
		"Help": &help,
	})
}
