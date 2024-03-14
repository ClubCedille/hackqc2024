package pages

import (
	"net/http"
	"slices"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
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
