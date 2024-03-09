package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func registerRoutes(r *gin.Engine, db *clover.DB) {

	r.GET("/", func(c *gin.Context) {
		dt := time.Now()
		// request := Request{
		// 	IP:       c.ClientIP(),
		// 	DateTime: dt.Format(time.RFC3339),
		// }
		// document := document.NewDocumentOf(request)
		// db.Insert(database.HackQcCollection, document)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Time": dt.Format("2006-01-02 15:04"),
		})

	})

	r.GET("/events-geojson", func(c *gin.Context) {
		if cachedGeoJSON == nil {
			fetchGeoJSON()
		}
		c.Data(http.StatusOK, "application/json", cachedGeoJSON)
	})

	r.GET("/map", func(c *gin.Context) {
		c.HTML(http.StatusOK, "map.html", nil)
	})

	r.Run()
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
