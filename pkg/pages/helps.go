package pages

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/comment"
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ClubCedille/hackqc2024/pkg/notifications"
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
	helpRequest := GetHelpFromContext(c, db)

	err := help.CreateHelp(db, helpRequest)
	if err != nil {
		log.Println("Error submitting help request:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	linkedEvent, err := event.GetEventById(db, helpRequest.EventId)
	if err != nil {
		log.Println("Error getting linked event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Help created successfully")
	notifMessage := fmt.Sprintf("Une offre aide a été soumise pour évènement près de vous %s.", helpRequest.MapObject.Name)
	notifications.NotifyEventSubscribers(
		db,
		notifMessage,
		linkedEvent.Subscribers,
	)

	c.Redirect(http.StatusSeeOther, "/map")
}

func GetHelpDetailAboutToBeModified(c *gin.Context, db *clover.DB) {
	id := c.Param("id")
	help, err := help.GetHelpById(db, id)
	if err != nil {
		log.Println("Error getting event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	// Format some data
	var formattedCoords []string
	for _, coord := range help.MapObject.Geometry.Coordinates {
		formattedCoords = append(formattedCoords, strconv.FormatFloat(coord, 'f', -1, 64))
	}
	fCoords := strings.Join(formattedCoords, ", ")

	help.MapObject.Description = strings.TrimSpace(help.MapObject.Description)

	c.HTML(http.StatusOK, "modals/update-help.html", gin.H{
		"Help":        help,
		"Coordinates": fCoords,
		"Tags":        strings.Join(help.MapObject.Tags, ", "),
	})
}

func UpdateHelp(c *gin.Context, db *clover.DB) {
	data := GetHelpFromContext(c, db)
	data.Id = c.Param("id")

	err := help.UpdateHelp(db, data)
	if err != nil {
		log.Println("Error updating help:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Help updated successfully")
}

func DeleteHelp(c *gin.Context, db *clover.DB) {
	helpId := c.Param("id")

	err := help.DeleteHelpById(db, helpId)
	if err != nil {
		log.Println("Error deleting help:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Help deleted successfully")
}

func GetHelpDetailsAboutToBeDelete(c *gin.Context, db *clover.DB) {
	id := c.Param("id")
	help, err := help.GetHelpById(db, id)
	if err != nil {
		log.Println("Error getting event:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "modals/delete-help.html", gin.H{
		"Help": &help,
	})
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
		"Helps":         helps,
		"ActiveSession": session.SessionIsActive(),
		"UserName":      session.ActiveSession.UserName,
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

	comments, err := comment.GetCommentsFormData(db, id)
	if err != nil {
		log.Println("Error fetching comments:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "modals/help-details.html", gin.H{
		"Help":       &help,
		"Comments":   comments,
		"IsLoggedIn": session.ActiveSession.AccountId != "",
	})
}

func GetHelpFromContext(c *gin.Context, db *clover.DB) help.Help {
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
			return help.Help{}
		}
		coordinatesArrayFloat = append(coordinatesArrayFloat, floatCoord)
	}

	contactInfos := c.PostForm("contact_infos")
	needHelp := c.PostForm("need_help") == "true"
	howToHelp := c.PostForm("how_to_help")
	howToUseHelp := c.PostForm("how_to_use_help")

	return help.Help{
		ContactInfos: contactInfos,
		NeedHelp:     needHelp,
		HowToHelp:    howToHelp,
		HowToUseHelp: howToUseHelp,
		EventId:      eventId,
		Modified:     false,
		MapObject: mapobject.MapObject{
			AccountId:   session.ActiveSession.AccountId,
			Name:        eventName,
			Description: eventDescription,
			Category:    eventCategory,
			Tags:        tagsArrayString,
			Geometry:    mapobject.Geometry{GeomType: "Point", Coordinates: coordinatesArrayFloat},
			Date:        time.Now(),
		},
	}
}

func PostCreateHelpComment(c *gin.Context, db *clover.DB) {

	err := comment.CreateComment(db, comment.Comment{
		Comment:  c.PostForm("comment"),
		OwnerId:  session.ActiveSession.AccountId,
		TargetId: c.PostForm("target_id"),
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Comment created successfully")

	comments, err := comment.GetCommentsFormData(db, c.PostForm("target_id"))
	if err != nil {
		log.Println("Error fetching comments:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	help, err := help.GetHelpById(db, c.PostForm("target_id"))
	if err != nil {
		log.Println("Error getting help:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.HTML(http.StatusOK, "modals/help-details.html", gin.H{
		"Help":       &help,
		"Comments":   comments,
		"IsLoggedIn": true,
	})
}

func ExportHelps(c *gin.Context, db *clover.DB) {
	idsStr := c.PostForm("ids")
	idsArray := strings.Split(idsStr, ",")

	var validIds []string
	for _, id := range idsArray {
		trimmedID := strings.TrimSpace(id)
		if trimmedID != "" {
			validIds = append(validIds, trimmedID)
		}
	}

	if len(validIds) == 0 {
		log.Println("No valid IDs provided for export")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid IDs provided for export"})
		return
	}

	SubmitHelpsToDC(c, db, validIds)

	log.Println("Helps exported successfully")
}
