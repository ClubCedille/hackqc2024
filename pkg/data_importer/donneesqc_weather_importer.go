package external_data

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
)

const DQC_WEATHER_NAME = "DonnéesQC Weather"

type DonneesQcWeatherEventSource struct{}

func (source DonneesQcWeatherEventSource) GetName() string {
	return DQC_WEATHER_NAME
}

func (source DonneesQcWeatherEventSource) GetAllEvents() ([]event.Event, error) {
	events := []event.Event{}

	weather, err := getWeatherData()

	if err != nil {
		return nil, err
	}

	for _, feature := range weather.Features {
		event, err := feature.ToEvent()
		if err != nil {
			msg := fmt.Sprintf("Error converting feature to event: %s", err.Error())
			fmt.Fprintln(os.Stderr, msg)
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

func getWeatherData() (*WeatherFeatureCollection, error) {
	params := map[string]string{
		"typeNames":    "msp_vigilance_crue_publique_v_type",
		"outputFormat": "geojson",
	}

	result, err := MakeWFSRequest("GetFeature", params)
	if err != nil {
		return nil, err
	}

	var weatherFeatures WeatherFeatureCollection
	err = json.Unmarshal(result, &weatherFeatures)
	if err != nil {
		return nil, err
	}

	return &weatherFeatures, err
}

type WeatherFeatureCollection struct {
	Name           string           `json:"name"`
	CollectionType string           `json:"type"`
	Features       []WeatherFeature `json:"features"`
}

type WeatherFeature struct {
	FeatureType string            `json:"type"`
	Properties  WeatherProperties `json:"properties"`
	Geometry    Geometry          `json:"geometry"`
}

func (feature *WeatherFeature) ToEvent() (event.Event, error) {
	date, err := time.Parse(DQC_TIME_FMT, feature.Properties.Date_mise_a_jour)
	if err != nil {
		return event.Event{}, err
	}

	return event.Event{
		DangerLevel: ParseSeverity(feature.Properties.Severite),
		UrgencyType: ParseUrgency(feature.Properties.Urgence),
		MapObject: mapobject.MapObject{
			Coordinates: FormatCoordinates(feature.Geometry.Coordinates),
			Polygon:     "",
			Name:        feature.Properties.Nom,
			Description: feature.Properties.Description,
			Category:    feature.Properties.Type,
			Tags:        []string{"external", DQC_WEATHER_NAME},
			Date:        date,
			AccountId:   SYSTEM_USER_GUID,
		},
	}, nil
}

type WeatherProperties struct {
	Nom              string `json:"nom"`
	Source           string `json:"source"`
	Territoire       string `json:"territoire"`
	Certitude        string `json:"certitude"`
	Severite         string `json:"severite"`
	Date_mise_a_jour string `json:"date_mise_a_jour"`
	Id_alerte        string `json:"id_alerte"`
	Urgence          string `json:"urgence"`
	Description      string `json:"description"`
	Type             string `json:"type"`
}
