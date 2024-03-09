package external_data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/event"
)

const VERSION = "1.1.0"
const BASE_URL = "https://geoegl.msp.gouv.qc.ca/ws/igo_gouvouvert.fcgi"

func getAllEvents() ([]event.Event, error) {
	events := []event.Event{}

	weather, err := getWeatherData()

	if err != nil {
		return nil, err
	}

	for _, feature := range weather.Features {
		events = append(events, feature.ToEvent())
	}

	return events, nil
}

func (feature *WeatherFeature) ToEvent() event.Event {
	return event.Event{
		DangerLevel: parseSeverity(feature.Properties.Severite),
		UrgencyType: parseUrgency(feature.Properties.Urgence),
	}
}

func parseUrgency(urgency string) event.UrgencyType {
	switch urgency {
	case "Future":
		return event.Futur
	case "Présent":
	case "Present":
		return event.Present
	}
	return event.Past
}

func parseSeverity(severity string) event.DangerLevel {
	switch severity {
	case "Importante":
		return event.High
	case "Modérée":
		return event.Medium
	}
	return event.Low
}

func getWeatherData() (*WeatherFeatureCollection, error) {
	params := map[string]string{
		"typeNames":    "msp_vigilance_crue_publique_v_type",
		"outputFormat": "geojson",
	}

	result, err := makeWFSRequest("GetFeature", params)
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

func makeWFSRequest(request string, params map[string]string) ([]byte, error) {
	params["version"] = VERSION
	params["service"] = "wfs"
	params["request"] = request

	queryString := toGetParams(params)
	request = BASE_URL + queryString
	resp, err := http.Get(request)

	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func toGetParams(params map[string]string) string {
	queryString := ""
	separator := "?"

	for key, value := range params {
		queryString += fmt.Sprintf("%s%s=%s", separator, key, value)
		separator = "&"
	}

	return queryString
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
}

type Geometry struct {
	GeomType    string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}
