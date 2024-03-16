package main

import (
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/pages"
	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func registerRoutes(r *gin.Engine, db *clover.DB) {
	r.Static("/static", "./templates/static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/map")
	})

	r.GET("/map", func(c *gin.Context) {
		pages.MapPage(c, db)
	})

	r.GET("/map-json", func(c *gin.Context) {
		pages.MapJson(c, db)
	})

	r.GET("/get-pannes-overlay", func(c *gin.Context) {
		pages.GetPannesOverlay(c, db)
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

	// r.GET("/submit-helps", func(c *gin.Context) {
	// 	pages.SubmitHelpsToDC(c, db)
	// })

	r.POST("/events/subscribe/:id", func(c *gin.Context) {
		pages.EventSubscribe(c, db)
	})

	r.GET("/event/:id", func(c *gin.Context) {
		pages.EventDetails(c, db)
	})

	r.GET("/help/:id", func(c *gin.Context) {
		pages.HelpDetails(c, db)
	})

	// The requests below require the user to be authenticated
	// Help
	r.POST("/create-help", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.CreateHelp(c, db)
	})

	// Event
	r.POST("/create-event", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.CreateEvent(c, db)
	})

	// Manage posts
	r.GET("/manage-post", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.GetManagedPost(c, db)
	})

	r.GET("/delete-event/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.GetEventDetailsAboutToBeDelete(c, db)
	})

	r.DELETE("/event/delete/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.DeleteEvent(c, db)
	})

	r.GET("/delete-help/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.GetHelpDetailsAboutToBeDelete(c, db)
	})

	r.DELETE("/help/delete/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.DeleteHelp(c, db)
	})

	r.GET("/update-event/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.GetEventDetailAboutToBeModified(c, db)
	})

	r.POST("/event/update/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.UpdateEvent(c, db)
	})

	r.GET("/update-help/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.GetHelpDetailAboutToBeModified(c, db)
	})

	r.POST("/help/update/:id", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.UpdateHelp(c, db)
	})

	r.POST("/export-helps", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.ExportHelps(c, db)
	})

	r.POST("/help/comment", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.PostCreateHelpComment(c, db)
	})

	r.POST("/event/comment", AuthRequiredMiddleware(), func(c *gin.Context) {
		pages.PostCreateEventComment(c, db)
	})
	r.GET("/a-propos", func(c *gin.Context) {
		c.HTML(http.StatusOK, "a-propos/index.html", gin.H{
			"ActiveSession": session.SessionIsActive(),
			"UserName":      session.ActiveSession.UserName,
		})
	})
}
