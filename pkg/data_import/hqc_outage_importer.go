package data_import

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
)

const HQC_URL = "http://pannes.hydroquebec.com/pannes/donnees/v3_0"
const HQC_TIME_FMT = "2006-01-02 15:04:05"
const HQC_OUTAGE_NAME = "Hydro-Québec Outage"

type HQCOutageEventSource struct{}

func (source HQCOutageEventSource) GetName() string {
	return HQC_OUTAGE_NAME
}

func (source HQCOutageEventSource) GetAllEvents() ([]event.Event, error) {
	date := time.Date(2024, 01, 01, 01, 01, 01, 01, time.UTC)
	return source.GetNewEventsFromDate(date)
}

// https://www.hydroquebec.com/documents-donnees/donnees-ouvertes/pannes-interruptions.html
func (source HQCOutageEventSource) GetNewEventsFromDate(date time.Time) ([]event.Event, error) {
	//Première requête : obtenir la version à jour du fichier BIS (« time stamp »)
	request := fmt.Sprintf("%s/bisversion.json", HQC_URL)
	resp, err := http.Get(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bisVersion := strings.Trim(string(body), "\"")

	//Deuxième requête : obtenir un fichier JSON contenant la liste des pannes.
	request = fmt.Sprintf("%s/bismarkers%s.json", HQC_URL, bisVersion)
	resp, err = http.Get(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var listePannes ListePannes
	err = json.Unmarshal(body, &listePannes)
	if err != nil {
		return nil, err
	}

	events := []event.Event{}
	for _, panne := range listePannes.Pannes {
		event, err := panne.ToEvent()
		if err != nil {
			msg := fmt.Sprintf("Error converting panne to event: %s", err.Error())
			fmt.Fprintln(os.Stderr, msg)
			continue
		}

		events = append(events, event)
	}

	return events, nil
}

type ListePannes struct {
	Pannes []Panne `json:"pannes"`
}

type Panne struct {
	NbAffectes    int
	DateDebut     time.Time
	DateFinPrevue time.Time
	CoordX        float64
	CoordY        float64
	Statut        string
	Cause         string
}

func (p *Panne) UnmarshalJSON(bs []byte) error {
	var err error

	arr := []interface{}{}
	json.Unmarshal(bs, &arr)

	p.NbAffectes = int(arr[0].(float64))

	dateDebutStr := arr[1].(string)
	if dateDebutStr == "" {
		p.DateDebut = time.Now()
	} else {
		p.DateDebut, err = time.Parse(HQC_TIME_FMT, arr[1].(string))
		if err != nil {
			return err
		}
	}

	dateFinStr := arr[2].(string)
	if dateFinStr == "" {
		p.DateFinPrevue = time.Time{}
	} else {
		p.DateFinPrevue, err = time.Parse(HQC_TIME_FMT, arr[2].(string))
		if err != nil {
			return err
		}
	}

	coordSplit := strings.Split(strings.Trim(arr[4].(string), "[]"), ",")
	p.CoordX, err = strconv.ParseFloat(coordSplit[0], 64)
	if err != nil {
		return err
	}

	p.CoordY, err = strconv.ParseFloat(strings.Trim(coordSplit[1], " "), 64)
	if err != nil {
		return err
	}

	codeStatus := arr[5].(string)
	switch codeStatus {
	case "A":
		p.Statut = "Travaux Assignés"
	case "L":
		p.Statut = "Équipe au travail"
	case "R":
		p.Statut = "Équipe en route"
	default:
		p.Statut = "Inconnu"
	}

	p.Cause = arr[6].(string)

	return nil
}

func (p Panne) ToEvent() (event.Event, error) {
	desc := fmt.Sprintf(
		"Clients affectés: %d\n Date de fin prévue: %s\n Statut: %s\nCause: %s",
		p.NbAffectes,
		p.DateFinPrevue.Format(HQC_TIME_FMT),
		p.Statut,
		p.Cause,
	)
	return event.Event{
		ExternalId: fmt.Sprintf("HYDROOUTAGE-%f_%f", p.CoordX, p.CoordY),
		MapObject: mapobject.MapObject{
			Geometry:    mapobject.Geometry{GeomType: "Point", Coordinates: []float64{p.CoordX, p.CoordY}},
			Name:        "Panne d'électricité",
			Description: desc,
			Category:    "Panne d'électricité",
			Tags:        []string{"external", HQC_OUTAGE_NAME},
			Date:        p.DateDebut,
			AccountId:   SYSTEM_USER_GUID,
		},
	}, nil
}
