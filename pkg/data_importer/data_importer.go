package external_data

import (
	"github.com/ClubCedille/hackqc2024/pkg/event"
)

const SYSTEM_USER_GUID = "dcf54266-549e-4546-9431-872f2d2fe8b3"

type EventSource interface {
	GetName() string
	GetAllEvents() ([]event.Event, error)

	//TODO: Incremental load with cron
	//GetNewEventsFromDate(date time.Time) ([]event.Event, error)
}

func GetEventSources() []EventSource {
	return []EventSource{
		DonneesQcWeatherEventSource{},
	}
}

func InitialLoad() ([]event.Event, error) {
	allEvents := []event.Event{}
	for _, source := range GetEventSources() {
		events, err := source.GetAllEvents()
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, events...)
	}
	return allEvents, nil
}
