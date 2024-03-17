package pages

import (
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func GetManagedPost(c *gin.Context, db *clover.DB) {
	var events []*event.Event
	events, err := event.GetAllEventsByAccountId(db, session.ActiveSession.AccountId)
	if err != nil {
		log.Println("Error fetching events:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	var helps []*help.Help
	helps, err = help.GetAllHelpsByAccountId(db, session.ActiveSession.AccountId)
	if err != nil {
		log.Println("Error fetching helps:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, e := range events {
		e.FlipCoords()
	}

	c.HTML(http.StatusOK, "profile/index.html", gin.H{
		"Events":        events,
		"Helps":         helps,
		"ActiveSession": session.SessionIsActive(),
		"UserName":      session.ActiveSession.UserName,
	})
}
