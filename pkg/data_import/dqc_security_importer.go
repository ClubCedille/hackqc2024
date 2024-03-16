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

const (
	DQC_SECURITY_NAME = "DonnéesQC Security"
	DQC_SECURITY_URL  = "http://geoegl.msp.gouv.qc.ca/apis/wss/historiquesc.fcgi"
	DQC_START_YEAR    = 2024
)

const MUN_POLYGONS_FILE = "switchedCoordinatesMunPolygons.json"

type DQCSecurityEventSource struct{}

func (source DQCSecurityEventSource) GetName() string {
	return DQC_SECURITY_NAME
}

func (source DQCSecurityEventSource) GetAllEvents() ([]event.Event, error) {
	date := time.Date(DQC_START_YEAR, 01, 01, 01, 01, 01, 01, time.UTC)
	events, err := source.GetNewEventsFromDate(date)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (source DQCSecurityEventSource) GetNewEventsFromDate(date time.Time) ([]event.Event, error) {
	filter := source.createDateFilterXml(date)
	params := map[string]string{
		"srsName":      "EPSG:4326",
		"typeNames":    "msp_risc_evenements_public",
		"outputFormat": "geojson",
		"filter":       filter,
	}
	securityData, err := searchSecurityData(params)
	events := parseSecurityData(securityData)

	if err != nil {
		return nil, err
	}

	return events, nil
}

func searchSecurityData(params map[string]string) (*SecurityFeatureCollection, error) {
	result, err := MakeWFSGetRequest(DQC_SECURITY_URL, "GetFeature", params)
	if err != nil {
		return nil, err
	}

	var securityData SecurityFeatureCollection
	err = json.Unmarshal(result, &securityData)
	if err != nil {
		return nil, err
	}

	return &securityData, err
}

func parseSecurityData(securityData *SecurityFeatureCollection) []event.Event {
	events := []event.Event{}

	municipalitiesPolygons := GetMunicipalityPolygons()

	for _, feature := range securityData.Features {
		event, err := feature.ToEvent()
		event.MunicipalityPolygon = municipalitiesPolygons[feature.Properties.Municipalite]
		if err != nil {
			msg := fmt.Sprintf("Error converting feature to event: %s", err.Error())
			fmt.Fprintln(os.Stderr, msg)
			continue
		}

		events = append(events, event)
	}

	return events
}

type SecurityFeatureCollection struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	Features []SecurityFeature `json:"features"`
}

type SecurityFeature struct {
	Id         int32              `json:"id"`
	Type       string             `json:"type"`
	Properties SecurityProperties `json:"properties"`
	Geometry   mapobject.Geometry `json:"geometry"`
}

type SecurityProperties struct {
	Code_alea               int32   `json:"code_alea"`
	Alea                    string  `json:"alea"`
	Code_municipalite       string  `json:"code_municipalite"`
	Municipalite            string  `json:"municipalite"`
	Precision_localisation  string  `json:"precision_localisation"`
	Info_compl_localisation string  `json:"info_compl_localisation"`
	Severite                string  `json:"severite"`
	Date_signalement        string  `json:"date_signalement"`
	Date_debut              string  `json:"date_debut"`
	Date_debut_imprecise    string  `json:"date_debut_imprecise"`
	Commentaire_date_debut  string  `json:"commentaire_date_debut"`
	Date_fin                string  `json:"date_fin"`
	Coord_x                 float32 `json:"coord_x"`
	Coord_y                 float32 `json:"coord_y"`
}

// calculate based on date?
func (f SecurityFeature) GetSeverity() event.UrgencyType {
	return event.Past
}

func GetMunicipalityPolygons() map[string][][][]float64 {
	municipalitiesPolygon, err := os.Open(MUN_POLYGONS_FILE)
	if err != nil {
		fmt.Println(err)
	}
	defer municipalitiesPolygon.Close()

	var municipalities map[string][][][]float64
	jsonParser := json.NewDecoder(municipalitiesPolygon)
	if err = jsonParser.Decode(&municipalities); err != nil {
		fmt.Println(err)
	}

	return municipalities
}

func (f SecurityFeature) GetDanger() event.DangerLevel {
	sev := strings.ToLower(f.Properties.Severite)
	if strings.Contains(sev, "important") || strings.Contains(sev, "extraordinaire") {
		return event.High
	} else if strings.Contains(sev, "possible") {
		return event.Medium
	} else {
		return event.Low
	}
}

func (f SecurityFeature) GetDate() (time.Time, error) {
	return time.Parse(DQC_DATE_FMT, f.Properties.Date_debut)
}

func (f SecurityFeature) ToEvent() (event.Event, error) {
	date, err := f.GetDate()

	if err != nil {
		return event.Event{}, err
	}

	return event.Event{
		DangerLevel: f.GetDanger(),
		UrgencyType: f.GetSeverity(),
		ExternalId:  fmt.Sprintf("Security-%d", f.Id),
		MapObject: mapobject.MapObject{
			Geometry:    f.Geometry,
			Name:        f.Properties.Alea,
			Description: "",
			Category:    f.Properties.Alea,
			Tags:        []string{"external", DQC_SECURITY_NAME},
			Date:        date,
			AccountId:   SYSTEM_USER_GUID,
		},
	}, nil
}

func (DQCSecurityEventSource) createDateFilterXml(date time.Time) string {
	searchParams := WFSSearch{
		LogicalOperator: "PropertyIsGreaterThan",
		PropertyName:    "date_signalement",
		PropertyValue:   date.Format(DQC_DATE_FMT),
		TypeName:        "msp_risc_evenements_public",
	}
	return searchParams.ToSecuritySearchPayload()
}

func (params WFSSearch) ToSecuritySearchPayload() string {
	return fmt.Sprintf(
		`<ogc:Filter xmlns:ogc="http://www.opengis.net/ogc" xmlns:gml="http://www.opengis.net/gml"><ogc:%s xmlns:ogc="http://www.opengis.net/ogc"><ogc:PropertyName xmlns:ogc="http://www.opengis.net/ogc">%s</ogc:PropertyName><ogc:Literal xmlns:ogc="http://www.opengis.net/ogc">%s</ogc:Literal></ogc:%s></ogc:Filter>`,
		params.LogicalOperator,
		params.PropertyName,
		params.PropertyValue,
		params.LogicalOperator)
}
