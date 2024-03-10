package main

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
	uuid "github.com/satori/go.uuid"

	"github.com/ClubCedille/hackqc2024/pkg/account"
	"github.com/ClubCedille/hackqc2024/pkg/data_import"
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
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

	// Initial load
	err = data_import.UpdateAll(db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	generateSeedData(db)

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"escapeSingleQuotes": func(s string) string {
			return strings.ReplaceAll(s, "'", "\\'")
		},
	})
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

func generateSeedData(db *clover.DB) {
	// Create an account
	err := account.CreateAccount(db, account.Account{
		Id:        uuid.NewV4().String(),
		UserName:  "sonoflope",
		FirstName: "son",
		LastName:  "oflope",
		Email:     "sonoflope@allo.com",
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

	accountExists, err := account.AccountExistsById(db, data_import.SYSTEM_USER_GUID)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if !accountExists {
		err = account.CreateAccount(db, account.Account{
			Id:        data_import.SYSTEM_USER_GUID,
			UserName:  "system",
			FirstName: "system",
			LastName:  "system",
			Email:     "system@allo.com",
		})
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Fetch the new account id
	docs, err := db.FindFirst(query.NewQuery(database.AccountCollection).Where(query.Field("user_name").Eq("sonoflope")))
	if err != nil {
		log.Fatalf(err.Error())
	}

	acc := account.Account{}
	docs.Unmarshal(&acc)

	// Create test help object
	err = help.CreateHelp(db, help.Help{
		Id: uuid.NewV4().String(),
		MapObject: mapobject.MapObject{
			AccountId:   acc.Id,
			Geometry:    mapobject.Geometry{GeomType: "Point", Coordinates: []float64{45.5017, -73.5673}},
			Name:        "Test help",
			Description: "This is a test help object",
			Category:    "Test",
			Tags:        []string{"test", "help"},
		},
		ContactInfos: "test contact infos",
		NeedHelp:     true,
		HowToHelp:    "test how to help",
		HowToUseHelp: "test how to use help",
	})
	if err != nil {
		log.Fatalf(err.Error())
	}

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
