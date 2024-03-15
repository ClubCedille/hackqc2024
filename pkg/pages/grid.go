package pages

import (
	"net/http"
	"slices"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
)

func GridPage(c *gin.Context, db *clover.DB) {

	// Get all events and helps
	docs, _ := db.FindAll(query.NewQuery(database.EventCollection).Sort(query.SortOption{Field: "map_object.date", Direction: -1}))
	events, _ := event.GetEventFromDocuments(docs)

	docs, _ = db.FindAll(query.NewQuery(database.HelpCollection).Sort(query.SortOption{Field: "map_object.date", Direction: -1}))
	helps, _ := help.GetHelpFromDocuments(docs)

	c.HTML(http.StatusOK, "grid/index.html", gin.H{
		"Events":        events,
		"Helps":         helps,
		"ActiveSession": session.SessionIsActive(),
		"UserName":      session.ActiveSession.UserName,
	})
}

func GridSearch(c *gin.Context, db *clover.DB) {

	// Search query

	var events []*event.Event
	var helps []*help.Help

	if slices.Contains(c.Request.URL.Query()["_.type"], "event") {
		events, _ = event.SearchEvents(db, c.Request.URL.Query(), false)
	}

	if slices.Contains(c.Request.URL.Query()["_.type"], "help") {
		helps, _ = help.SearchHelps(db, c.Request.URL.Query(), false)
	}

	c.HTML(http.StatusOK, "components/grid", gin.H{
		"Events": events,
		"Helps":  helps,
	})
}
