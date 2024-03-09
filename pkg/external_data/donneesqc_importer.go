package external_data

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/event"
)

const VERSION = "1.1.0"
const BASE_URL = "https://geoegl.msp.gouv.qc.ca/ws/igo_gouvouvert.fcgi"

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

func (feature *WeatherFeature) ToEvent() event.Event {
	return event.Event{
		DangerLevel: ParseSeverity(feature.Properties.Severite),
		UrgencyType: ParseUrgency(feature.Properties.Urgence),
	}
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

func GetAllExternalEvents() ([]event.Event, error) {
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

	queryString := ToGetParams(params)
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
