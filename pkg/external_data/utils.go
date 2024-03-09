package external_data

import (
	"fmt"

	"github.com/ClubCedille/hackqc2024/pkg/event"
)

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
