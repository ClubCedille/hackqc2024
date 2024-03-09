package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	uuid "github.com/satori/go.uuid"

	"github.com/ClubCedille/hackqc2024/pkg/account"
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/event"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ostafen/clover/v2/query"
)

type Request struct {
	IP       string
	DateTime string
}

func main() {
	// Init database
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Create an account
	err = account.CreateAccount(db, account.Account{
		Id:        uuid.NewV4().String(),
		UserName:  "sonoflope",
		FirstName: "son",
		LastName:  "oflope",
		Email:     "sonoflope@allo.com",
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Fetch the new account id
	docs, err := db.FindFirst(query.NewQuery(database.AccountCollection).Where(query.Field("user_name").Eq("sonoflope")))
	if err != nil {
		log.Fatalf(err.Error())
	}

	acc := account.Account{}
	docs.Unmarshal(&acc)

	// Create map object
	mapOjb := mapobject.MapObject{
		Coordinates: "test",
		Polygon:     "test",
		Name:        "this is a test",
		Description: "this is a test",
		Category:    "this is a test",
		Tags:        []string{"test1", "test2"},
		Date:        time.Now(),
		AccountId:   acc.Id,
	}
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Create an event
	err = event.CreateEvent(db, event.Event{
		Id:          uuid.NewV4().String(),
		DangerLevel: event.DangerLevel(1),
		UrgencyType: event.UrgencyType(1),
		MapObject:   mapOjb,
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/**/*.html")

	registerRoutes(r, db)

	err = r.Run()
	if err != nil {
		fmt.Print("Failed to run")
		return
	}

	defer func(db *clover.DB) {
		err := db.Close()
		if err != nil {
			fmt.Print("Failed to close CloverDB during program exit!")
			return
		}
	}(db)
}

// Temp example of fetching from données Québec
// var cachedGeoJSON []byte

// func fetchGeoJSON() {
// 	if cachedGeoJSON == nil {
// 		resp, err := http.Get("https://donnees.montreal.ca/dataset/6a4cbf2c-c9b7-413a-86b1-e8f7081e2578/resource/35307457-a00f-4912-9941-8095ead51f6e/download/evenements.geojson")
// 		if err != nil {
// 			log.Println("Error fetching GeoJSON:", err)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		data, err := io.ReadAll(resp.Body)
// 		if err != nil {
// 			log.Println("Error reading GeoJSON:", err)
// 			return
// 		}
// 		cachedGeoJSON = data
// 	}
// }
