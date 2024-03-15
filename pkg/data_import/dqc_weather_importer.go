package data_import

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	jsoniter "github.com/json-iterator/go"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/geojson"
)

const (
	DQC_WEATHER_NAME = "DonnéesQC Weather"
	DQC_WEATHER_URL  = "https://geoegl.msp.gouv.qc.ca/ws/igo_gouvouvert.fcgi"
)

var config = jsoniter.Config{
	EscapeHTML:              true,
	SortMapKeys:             false,
	MarshalFloatWith6Digits: true,
}.Froze()

type DQCWeatherEventSource struct{}

func (source DQCWeatherEventSource) GetName() string {
	return DQC_WEATHER_NAME
}

func (source DQCWeatherEventSource) GetAllEvents() ([]event.Event, error) {
	events := []event.Event{}

	weather, err := getWeatherData(map[string]string{})
	alerts, err := getWeatherAlerts(map[string]string{})

	if err != nil {
		return nil, err
	}

	for _, weather_feature := range weather.Features {

		// Find the alert for this feature by alert_id
		alert_id := strings.Split(weather_feature.Properties.Id_alerte, ".")[1]
		var foundAlert *geojson.Feature
		for _, alert := range alerts.Features {
			if strings.Contains(alert.Properties["id_alerte"].(string), alert_id) {
				foundAlert = alert
				break
			}
		}

		// hack: convert coordinate to wrong float64 array for db
		// see: https://github.com/paulmach/orb/issues/45
		var coordinates []float64
		// Replace the geometry of the feature (single point) with the alert geometry (polygon)
		if foundAlert != nil {
			geometry := foundAlert.Geometry
			for _, ring := range geometry.(orb.Polygon) {
				for _, point := range ring {
					coordinates = append(coordinates, point.X(), point.Y()) //massacre of geojson
				}
			}
			weather_feature.Geometry.GeomType = "Polygon"
			weather_feature.Geometry.Coordinates = coordinates
		}

		event, err := weather_feature.ToEvent()
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

/*
These are the WFS query params that are used to get weather alerts with polygons from the DQC WFS service.
SERVICE: WFS
REQUEST: GetFeature
VERSION: 2.0.0
TYPENAMES: ms:masas_naad_adna_s_public
SRSNAME: urn:ogc:def:crs:EPSG::4326
BBOX: -6050013.3284445209428668,-1035657.34829420340247452,4471994.56733370572328568,5149718.85507926531136036,urn:ogc:def:crs:EPSG::32198
outputFormat: geojson
*/

// This gets the weather alerts from layer "ms:masas_naad_adna_s_public" on the DQC WFS service.
func getWeatherAlerts(params map[string]string) (*geojson.FeatureCollection, error) {
	params["typeNames"] = "ms:masas_naad_adna_s_public"
	params["outputFormat"] = "geojson"
	params["srsName"] = "EPSG:4326"
	params["bbox"] = "-6050013.3284445209428668,-1035657.34829420340247452,4471994.56733370572328568,5149718.85507926531136036,urn:ogc:def:crs:EPSG::32198"

	result, err := MakeWFSGetRequest(DQC_WEATHER_URL, "GetFeature", params)
	if err != nil {
		return nil, err
	}

	fc, _ := geojson.UnmarshalFeatureCollection(result)

	return fc, err
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
