package pages

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/help"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type GeoJSONPair struct {
	GeoJson GeoJSON `json:"geoJson"`
	Style   Style   `json:"style"`
}

type Style struct {
	Color    string `json:"color"`
	Icon     string `json:"icon"`
	IconSize int    `json:"iconSize"` //0-3
}

type GeoJSON struct {
	Type       string              `json:"type"`
	Geometry   mapobject.Geometry  `json:"geometry"`
	Properties mapobject.MapObject `json:"properties"`
}

type NameValue struct {
	Name  string
	Value string
}

// Styling with google material icons
// using list at all_material_icons.txt
var CategoryStyles = map[string]Style{
	"Vent": {
		Color:    "blue",
		IconSize: 0,
		Icon:     "air",
	},
	"Pluie": {
		Color:    "blue",
		IconSize: 0,
		Icon:     "rainy",
	},
	"Neige": {
		Color:    "blue",
		IconSize: 0,
		Icon:     "ac_unit",
	},
	"Tempête hivernale": {
		Color:    "red",
		IconSize: 2,
		Icon:     "severe_cold",
	},
	"Onde de tempête": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "weather_mix",
	},
	"Inondation": {
		Color:    "red",
		IconSize: 2,
		Icon:     "flood",
	},
	"Panne d'électricité": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "electric_bolt",
	},
	"Mouvement de terrain": {
		Color:    "red",
		IconSize: 2,
		Icon:     "landslide",
	},
	"Autre": {
		Color:    "red",
		IconSize: 1,
		Icon:     "warning",
	},
	"Orage violent": {
		Color:    "red",
		IconSize: 2,
		Icon:     "thunderstorm",
	},
	"Fermeture de route": {
		Color:    "yellow",
		IconSize: 0,
		Icon:     "traffic",
	},
	"Matières dangereuses": {
		Color:    "red",
		IconSize: 2,
		Icon:     "skull",
	},
	"Panne de télécommunication": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "wifi",
	},
	"Feu de forêt": {
		Color:    "red",
		IconSize: 2,
		Icon:     "local_fire_department",
	},
	"Érosion": {
		Color:    "yellow",
		IconSize: 0,
		Icon:     "landslide",
	},
	"Tremblement de terre": {
		Color:    "red",
		IconSize: 2,
		Icon:     "earthquake",
	},
	"Alimentation publique en eau potable": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "water",
	},
	"Incendie industriel": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "local_fire_department",
	},
	"Pluie verglaçante": {
		Color:    "blue",
		IconSize: 0,
		Icon:     "weather_mix",
	},
	"Fermeture de pont": {
		Color:    "yellow",
		IconSize: 0,
		Icon:     "traffic",
	},
	"Accident de voiture": {
		Color:    "yellow",
		IconSize: 0,
		Icon:     "traffic",
	},
	"Écrasement d'avion": {
		Color:    "yellow",
		IconSize: 0,
		Icon:     "flight_land",
	},
	"Risque d'explosion": {
		Color:    "red",
		IconSize: 2,
		Icon:     "destruction",
	},
	"Débordement de barrage": {
		Color:    "red",
		IconSize: 2,
		Icon:     "flood",
	},
	"Feu urbain": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "local_fire_department",
	},
	"Accident maritime": {
		Color:    "yellow",
		IconSize: 0,
		Icon:     "traffic",
	},
	"Effondrement de structure": {
		Color:    "red",
		IconSize: 2,
		Icon:     "traffic",
	},
	"Tornade": {
		Color:    "red",
		IconSize: 2,
		Icon:     "tornado",
	},
	"Crise civile": {
		Color:    "red",
		IconSize: 2,
		Icon:     "emergency_home",
	},
	"Risque de gaz toxiques": {
		Color:    "red",
		IconSize: 2,
		Icon:     "emergency_home",
	},
	"Accident ferroviaire": {
		Color:    "yellow",
		IconSize: 0,
		Icon:     "traffic",
	},
	"Ouragan": {
		Color:    "red",
		IconSize: 2,
		Icon:     "cyclone",
	},
	"Qualité de l'air": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "airwave",
	},
	"Santé": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "medical_services",
	},
	"Vague de chaleur": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "sunny",
	},
	"Vague de froid": {
		Color:    "yellow",
		IconSize: 1,
		Icon:     "ac_unit",
	},
	"Avalanche": {
		Color:    "red",
		IconSize: 2,
		Icon:     "landslide",
	},
	"Maladie infectieuse": {
		Color:    "red",
		IconSize: 2,
		Icon:     "coronavirus",
	},
}

func sortWithAccents(s []string) {
	// Create a Collator for French
	fr := collate.New(language.French, collate.Loose)

	// Sort the strings
	fr.SortStrings(s)
}

func MapPage(c *gin.Context, db *clover.DB) {
	filters := c.Request.URL.Query()
	mapItems, err := retrieveMapItems(db, filters)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	jsonValue, err := json.Marshal(mapItems)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	categoryKeys := make([]string, 0, len(CategoryStyles))
	for k := range CategoryStyles {
		categoryKeys = append(categoryKeys, k)
	}
	sortWithAccents(categoryKeys)

	// For create event form
	mapCategories := make([]interface{}, len(categoryKeys))
	for i, key := range categoryKeys {
		category := NameValue{
			Name:  key,
			Value: key,
		}
		mapCategories[i] = category
	}

	fmt.Println("mapCategories", mapCategories)

	urgencyLevels := []NameValue{
		{
			Name:  "Futur",
			Value: fmt.Sprint(event.Futur),
		},
		{
			Name:  "Passé",
			Value: fmt.Sprint(event.Past),
		},
		{
			Name:  "Présent",
			Value: fmt.Sprint(event.Present),
		},
	}

	dangerLevels := []NameValue{
		{
			Name:  "Élevé",
			Value: fmt.Sprint(event.High),
		},
		{
			Name:  "Modéré",
			Value: fmt.Sprint(event.Medium),
		},
		{
			Name:  "Faible",
			Value: fmt.Sprint(event.Low),
		},
	}

	c.HTML(http.StatusOK, "map/index.html", gin.H{
		"MapItemsJson":  string(jsonValue),
		"Categories":    mapCategories,
		"ActiveSession": session.ActiveSession.UserName,
		"MapCategory":   mapCategories,
		"UrgencyLevels": urgencyLevels,
		"DangerLevels":  dangerLevels,
		"CategoryKeys":  categoryKeys, // for selection in event form
	})
}

func MapJson(c *gin.Context, db *clover.DB) {
	filters := c.Request.URL.Query()
	mapItems, err := retrieveMapItems(db, filters)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, mapItems)
}

func retrieveMapItems(db *clover.DB, filters map[string][]string) ([]GeoJSONPair, error) {
	var events []*event.Event
	var helps []*help.Help
	var err error

	if filters["type"] == nil || filters["type"][0] == "events" {
		events, err = event.GetEventWithFilters(db, filters, true)
		if err != nil {
			return nil, err
		}
	}

	if filters["type"] == nil || filters["type"][0] == "helps" {
		helps, err = help.GetHelpWithFilters(db, filters, true)
		if err != nil {
			return nil, err
		}
	}

	evSize := len(events)
	helpSize := len(helps)
	mapItems := make([]GeoJSONPair, evSize+helpSize)

	for i, v := range events {
		styleEvent, exists := CategoryStyles[v.MapObject.Category]
		if !exists {
			styleEvent = Style{
				Color:    "red",
				Icon:     "location_on",
				IconSize: 1,
			}
		}

		mapItems[i] = GeoJSONPair{
			GeoJson: GeoJSON{
				Type:       "Feature",
				Geometry:   v.MapObject.Geometry,
				Properties: v.MapObject,
			},
			Style: styleEvent,
		}
	}

	for i, v := range helps {
		styleHelp, exists := CategoryStyles[v.MapObject.Category]
		if !exists {
			styleHelp = Style{
				Color:    "green",
				Icon:     "location_on",
				IconSize: 1,
			}
		}
		mapItems[i+evSize] = GeoJSONPair{
			GeoJson: GeoJSON{
				Type:       "Feature",
				Geometry:   v.MapObject.Geometry,
				Properties: v.MapObject,
			},
			Style: styleHelp,
		}
	}

	return mapItems, nil
}

func GetPannesOverlay(c *gin.Context, db *clover.DB) {
	files, err := os.ReadDir("tmp/")
	if err != nil {
		log.Println("Error reading directory:", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	matchingFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), "outageAreas") {
			matchingFiles = append(matchingFiles, file.Name())
		}
	}
	sort.Slice(matchingFiles, func(i, j int) bool {
		return matchingFiles[i] < matchingFiles[j]
	})
	if len(matchingFiles) > 0 {
		c.File("tmp/" + matchingFiles[0])
	} else {
		log.Println("No matching files found")
		c.Status(http.StatusInternalServerError)
	}
}
