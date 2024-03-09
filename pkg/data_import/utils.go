package data_import

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/event"
)

const DQC_VERSION = "1.1.0"
const DQC_BASE_URL = "https://geoegl.msp.gouv.qc.ca/ws/igo_gouvouvert.fcgi"
const DQC_TIME_FMT = "2006/01/02 15:04"

func ToGetParams(params map[string]string) string {
	queryString := ""
	separator := "?"

	for key, value := range params {
		queryString += fmt.Sprintf("%s%s=%s", separator, key, value)
		separator = "&"
	}

	return queryString
}

func ParseUrgency(urgency string) event.UrgencyType {
	switch urgency {
	case "Future":
		return event.Futur
	case "Présent":
	case "Present":
		return event.Present
	}
	return event.Past
}

func ParseSeverity(severity string) event.DangerLevel {
	switch severity {
	case "Importante":
		return event.High
	case "Modérée":
		return event.Medium
	}
	return event.Low
}

func MakeWFSRequest(request string, params map[string]string) ([]byte, error) {
	params["version"] = DQC_VERSION
	params["service"] = "wfs"
	params["request"] = request

	queryString := ToGetParams(params)
	request = DQC_BASE_URL + queryString
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

type Geometry struct {
	GeomType    string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

func FormatCoordinates(coordinates []float64) string {
	return fmt.Sprintf("%f,%f", coordinates[0], coordinates[1])
}
