package data_import

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ClubCedille/hackqc2024/pkg/event"
)

const DQC_VERSION = "2.0.0"
const DQC_TIME_FMT = "2006/01/02 15:04"
const DQC_DATE_FMT = "2006-01-02"

func ToGetParams(params map[string]string) string {
	queryString := ""
	separator := "?"

	for key, value := range params {
		value = url.QueryEscape(value)
		value = strings.ReplaceAll(value, "+", "%20")
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

func MakeWFSGetRequest(baseUrl string, request string, params map[string]string) ([]byte, error) {
	params["version"] = DQC_VERSION
	params["service"] = "wfs"
	params["request"] = request

	queryString := ToGetParams(params)
	request = baseUrl + queryString
	client := &http.Client{}
	req, err := http.NewRequest("GET", request, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "definitely-not-Go-htt-p-client/1.1")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)

		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, err
}

func MakeWFSPostRequest(baseUrl string, body string, params map[string]string) ([]byte, error) {
	params["version"] = DQC_VERSION

	queryString := ToGetParams(params)
	request := baseUrl + queryString
	resp, err := http.Post(request, "application/xml", bytes.NewBuffer([]byte(body)))

	if err != nil || resp.StatusCode != 200 {
		return nil, err
	}

	defer resp.Body.Close()
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return resBody, err
}

type Geometry struct {
	GeomType    string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type WFSSearch struct {
	LogicalOperator string
	PropertyName    string
	PropertyValue   string
	TypeName        string
}

func DownloadFile(url string, filepath string) error {
	zip_out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer zip_out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(zip_out, resp.Body)
	return err
}
