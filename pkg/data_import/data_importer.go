package data_import

import (
	"fmt"
	"os"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/event"
	"github.com/ClubCedille/hackqc2024/pkg/watermark"
	"github.com/ostafen/clover/v2"
)

const SYSTEM_USER_GUID = "dcf54266-549e-4546-9431-872f2d2fe8b3"

type EventSource interface {
	GetName() string
	GetAllEvents() ([]event.Event, error)

	//TODO: Incremental load with cron
	GetNewEventsFromDate(date time.Time) ([]event.Event, error)
}

func GetEventSources() []EventSource {
	return []EventSource{
		// DQCWeatherEventSource{},
		// DQCSecurityEventSource{},
		HQCOutageEventSource{},
	}
}

func GetEventSourceByName(name string) EventSource {
	for _, source := range GetEventSources() {
		if source.GetName() == name {
			return source
		}
	}

	return nil
}

func UpdateAll(db *clover.DB) error {
	for _, source := range GetEventSources() {
		events, err := UpdateFromSource(db, source)
		if err != nil {
			msg := fmt.Sprintf("Error updating events from source %s: %s", source.GetName(), err.Error())
			fmt.Fprintln(os.Stderr, msg)
			continue
		}

		for _, e := range events {
			err = createOrUpdateExternalEvent(db, e)
			if err != nil {
				msg := fmt.Sprintf("Error creating or updating external event %s: %s", e.ExternalId, err.Error())
				fmt.Fprintln(os.Stderr, msg)
				continue
			}
		}
	}

	return nil
}

func UpdateFromSource(db *clover.DB, source EventSource) ([]event.Event, error) {
	watermarkExists, err := watermark.WatermarkExistsByName(db, source.GetName())
	if err != nil {
		return nil, err
	}

	var events []event.Event
	var watermarkObj watermark.Watermark
	if watermarkExists {
		watermarkObj, err = watermark.GetWatermark(db, source.GetName())
		if err != nil {
			return nil, err
		}
		events, err = source.GetNewEventsFromDate(watermarkObj.Watermark)
		if err != nil {
			return nil, err
		}
	} else {
		events, err = source.GetAllEvents()
		if err != nil {
			return nil, err
		}

		watermarkObj = watermark.Watermark{
			Name: source.GetName(),
		}
	}

	watermarkObj.Watermark = time.Now()
	if watermarkObj.Id == "" {
		err = watermark.CreateWatermark(db, watermarkObj)
	} else {
		err = watermark.UpdateWatermark(db, watermarkObj)
	}
	if err != nil {
		return nil, err
	}

	return events, nil
}

func createOrUpdateExternalEvent(db *clover.DB, updatedEvent event.Event) error {
	eventExists, err := event.EventExistsByExternalId(db, updatedEvent.ExternalId)

	if err != nil {
		return err
	}

	//Not sure if legit, but err means that the event doesn't exist (I think)
	if eventExists {
		existingEvent, err := event.GetEventByExternalId(db, updatedEvent.ExternalId)

		if err != nil {
			return err
		}

		updatedEvent.Id = existingEvent.Id
		return event.UpdateEvent(db, updatedEvent)
	} else {
		return event.CreateEvent(db, updatedEvent)
	}
}
