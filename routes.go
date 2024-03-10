package main

import (
	"io"
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/pages"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func registerRoutes(r *gin.Engine, db *clover.DB) {
	r.Static("/static", "./templates/static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/map")
	})

	r.GET("/events-geojson", func(c *gin.Context) {
		if cachedGeoJSON == nil {
			fetchGeoJSON()
		}
		c.Data(http.StatusOK, "application/json", cachedGeoJSON)
	})

	r.GET("/map", func(c *gin.Context) {
		pages.MapPage(c, db)
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

	r.POST("/create-event", func(c *gin.Context) {
		pages.CreateEvent(c, db)
	})

	r.POST("/update-event", func(c *gin.Context) {
		pages.UpdateEvent(c, db)
	})

	r.DELETE("/delete-event", func(c *gin.Context) {
		pages.DeleteEvent(c, db)
	})

	// Help
	r.GET("/helps", func(c *gin.Context) {
		pages.HelpPage(c, db)
	})

	r.POST("/create-help", func(c *gin.Context) {
		pages.CreateHelp(c, db)
	})

	r.POST("/update-help", func(c *gin.Context) {
		pages.UpdateHelp(c, db)
	})

	r.DELETE("/delete-help", func(c *gin.Context) {
		pages.DeleteHelp(c, db)
	})

	// Account
	r.POST("/create-account", func(c *gin.Context) {
		pages.CreateAccount(c, db)
	})

	r.POST("/update-account", func(c *gin.Context) {
		pages.UpdateAccount(c, db)
	})

}

// Temp example of fetching from données Québec
var cachedGeoJSON []byte

func fetchGeoJSON() {

	if cachedGeoJSON == nil {
		resp, err := http.Get("https://donnees.montreal.ca/dataset/6a4cbf2c-c9b7-413a-86b1-e8f7081e2578/resource/35307457-a00f-4912-9941-8095ead51f6e/download/evenements.geojson")
		if err != nil {
			log.Println("Error fetching GeoJSON:", err)
			return
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading GeoJSON:", err)
			return
		}
		cachedGeoJSON = data
	}
}
