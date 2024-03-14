package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/query"
	uuid "github.com/satori/go.uuid"

	"github.com/ClubCedille/hackqc2024/pkg/account"
	"github.com/ClubCedille/hackqc2024/pkg/data_import"
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
)

const (
	GIN_SESSION_NAME   = "gin-session"
	GIN_SESSION_SECRET = "gin-session-secret"
)

type Request struct {
	IP       string
	DateTime string
}

const TMP_DIR = "tmp"

func main() {
	// Init database
	db, err := database.InitDatabase()
	if err != nil {
		log.Fatalf(err.Error())
	}

	os.Mkdir("tmp", 0777)
	// Initial load
	err = data_import.UpdateAll(db)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// generateSeedData(db)

	r := gin.Default()
	r.SetFuncMap(template.FuncMap{
		"escapeSingleQuotes": func(s string) string {
			return strings.ReplaceAll(s, "'", "\\'")
		},
	})

	store := cookie.NewStore([]byte(GIN_SESSION_SECRET))
	store.Options(sessions.Options{MaxAge: 60 * 60 * 24}) // expire in a day

	r.Use(sessions.Sessions(GIN_SESSION_NAME, store))

	r.Use(LoginMiddleware())

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
	// testHelp, err := db.FindFirst(query.NewQuery(database.HelpCollection).Where(query.Field("map_object.name").Eq("Test help")))
	// if err == nil && testHelp == nil {
	// 	err = help.CreateHelp(db, help.Help{
	// 		Id: uuid.NewV4().String(),
	// 		MapObject: mapobject.MapObject{
	// 			AccountId:   acc.Id,
	// 			Geometry:    mapobject.Geometry{GeomType: "Point", Coordinates: []float64{-73.5673, 45.5017}},
	// 			Name:        "Test - Utilisez ma maison",
	// 			Description: "J'ai deux chambres à coucher de libre",
	// 			Category:    "Hébergement",
	// 			Tags:        []string{"test", "help"},
	// 		},
	// 		ContactInfos: "555 444 3333",
	// 		NeedHelp:     true,
	// 		HowToHelp:    "N/A",
	// 		HowToUseHelp: "Venez chez moi.",
	// 	})
	// }
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
}
