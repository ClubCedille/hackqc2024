package data_import

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
)

const DQC_WEATHER_NAME = "DonnéesQC Weather"
const DQC_WEATHER_URL = "https://geoegl.msp.gouv.qc.ca/ws/igo_gouvouvert.fcgi"

type DQCWeatherEventSource struct{}

func (source DQCWeatherEventSource) GetName() string {
	return DQC_WEATHER_NAME
}

func (source DQCWeatherEventSource) GetAllEvents() ([]event.Event, error) {
	events := []event.Event{}

	weather, err := getWeatherData(map[string]string{})

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

func (source DQCWeatherEventSource) GetNewEventsFromDate(date time.Time) ([]event.Event, error) {
	body := source.createDateFilterXml(date)
	weather, err := searchWeatherData(body)
	if err != nil {
		return nil, err
	}

	events := parseWeather(weather)

	return events, nil
}

func parseWeather(weather *WeatherFeatureCollection) []event.Event {
	events := []event.Event{}

	for _, feature := range weather.Features {
		event, err := feature.ToEvent()
		if err != nil {
			msg := fmt.Sprintf("Error converting feature to event: %s", err.Error())
			fmt.Fprintln(os.Stderr, msg)
			continue
		}

		events = append(events, event)
	}

	return events
}

func searchWeatherData(body string) (*WeatherFeatureCollection, error) {
	result, err := MakeWFSPostRequest(DQC_WEATHER_URL, body, map[string]string{})
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

func getWeatherData(params map[string]string) (*WeatherFeatureCollection, error) {
	params["typeNames"] = "msp_vigilance_crue_publique_v_type"
	params["outputFormat"] = "geojson"
	params["srsName"] = "EPSG:4326"

	result, err := MakeWFSGetRequest(DQC_WEATHER_URL, "GetFeature", params)
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
	FeatureType string             `json:"type"`
	Properties  WeatherProperties  `json:"properties"`
	Geometry    mapobject.Geometry `json:"geometry"`
}

func (feature *WeatherFeature) ToEvent() (event.Event, error) {
	date, err := time.Parse(DQC_TIME_FMT, feature.Properties.Date_mise_a_jour)
	if err != nil {
		return event.Event{}, err
	}

	return event.Event{
		DangerLevel: ParseSeverity(feature.Properties.Severite),
		UrgencyType: ParseUrgency(feature.Properties.Urgence),
		ExternalId:  feature.Properties.Id_alerte,
		MapObject: mapobject.MapObject{
			Geometry:    feature.Geometry,
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

func (DQCWeatherEventSource) createDateFilterXml(date time.Time) string {
	searchParams := WFSSearch{
		LogicalOperator: "PropertyIsGreaterThan",
		PropertyName:    "date_mise_a_jour",
		PropertyValue:   date.Format(DQC_TIME_FMT),
		TypeName:        "msp_vigilance_crue_publique_v_type",
	}
	return strings.Trim(searchParams.ToWeatherSearchPayload("EPSG:4326", "geojson"), "\n ")
}

func (params WFSSearch) ToWeatherSearchPayload(srsName string, outputFormat string) string {
	return fmt.Sprintf(
		`<?xml version="1.0"?>
	<wfs:GetFeature xmlns:wfs="http://www.opengis.net/wfs/2.0" xmlns:fes="http://www.opengis.net/fes/2.0" xmlns:gml="http://www.opengis.net/gml/3.2" xmlns:sf="http://www.openplans.org/spearfish" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.opengis.net/wfs/2.0 http://schemas.opengis.net/wfs/2.0/wfs.xsd         http://www.opengis.net/gml/3.2 http://schemas.opengis.net/gml/3.2.1/gml.xsd" service="WFS" version="%s" outputFormat="%s">
	  <wfs:Query typeNames="%s" srsName="%s">
		<fes:Filter>
		  <%s>
			<ValueReference>%s</ValueReference>
			<Literal>%s</Literal>
		  </%s>
		</fes:Filter>
	  </wfs:Query>
	</wfs:GetFeature>`,
		DQC_VERSION,
		outputFormat,
		params.TypeName,
		srsName,
		params.LogicalOperator,
		params.PropertyName,
		params.PropertyValue,
		params.LogicalOperator,
	)
}
