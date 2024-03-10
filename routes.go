package main

import (
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/pages"
	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func authRegisterRoutes(r *gin.Engine, group *gin.RouterGroup, db *clover.DB) {
	group.Use(AuthRequiredMiddleware())
	{
		// Help
		group.POST("/create-help", func(c *gin.Context) {
			pages.CreateHelp(c, db)
		})

		group.POST("/update-help", func(c *gin.Context) {
			pages.UpdateHelp(c, db)
		})

		group.DELETE("/delete-help", func(c *gin.Context) {
			pages.DeleteHelp(c, db)
		})

		// Event
		group.POST("/create-event", func(c *gin.Context) {
			pages.CreateEvent(c, db)
		})

		group.POST("/update-event", func(c *gin.Context) {
			pages.UpdateEvent(c, db)
		})

		group.DELETE("/delete-event", func(c *gin.Context) {
			pages.DeleteEvent(c, db)
		})
	}
}

func registerRoutes(r *gin.Engine, db *clover.DB) {
	r.Static("/static", "./templates/static")

	r.GET("/", func(c *gin.Context) {
		session.GetActiveSession(c)
		c.Redirect(http.StatusSeeOther, "/map")
	})

	r.GET("/map", func(c *gin.Context) {
		pages.MapPage(c, db)
	})

	r.GET("/map-json", func(c *gin.Context) {
		pages.MapJson(c, db)
	})

	// Event-Help Grid
	r.GET("/grid", func(c *gin.Context) {
		pages.GridPage(c, db)
	})

	r.GET("/grid/search", func(c *gin.Context) {
		pages.GridSearch(c, db)
	})

	// Events
	r.GET("/events/table", func(c *gin.Context) {
		pages.EventTablePage(c, db)
	})

	r.GET("/events/table/search", func(c *gin.Context) {
		pages.SearchEventTable(c, db)
	})

	r.GET("/eventCards", func(c *gin.Context) {
		pages.EventsPage(c, db)
	})

	r.GET("/create-event", func(c *gin.Context) {
		pages.GetCreateEvent(c, db)
	})

	// Account
	r.GET("/create-account", func(c *gin.Context) {
		pages.GetCreateAccount(c)
	})

	r.POST("/create-account", func(c *gin.Context) {
		pages.CreateAccount(c, db)
	})

	r.POST("/update-account", func(c *gin.Context) {
		pages.UpdateAccount(c, db)
	})

	r.GET("/login", func(c *gin.Context) {
		pages.GetLogin(c)
	})

	r.POST("/login", func(c *gin.Context) {
		pages.Login(c, db)
	})

	r.POST("/logout", func(c *gin.Context) {
		pages.Logout(c)
	})

	// Help
	r.GET("/helps", func(c *gin.Context) {
		pages.HelpPage(c, db)
	})

	r.GET("/helps/table", func(c *gin.Context) {
		pages.HelpTablePage(c, db)
	})

	r.GET("/submit-events", func(c *gin.Context) {
		pages.SubmitEvents(c, db)
	})
}
