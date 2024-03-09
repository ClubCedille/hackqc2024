package external_data

import (
	"fmt"
	"testing"
)

func TestToGetParams(t *testing.T) {
	params := map[string]string{
		"test1": "val1",
		"test2": "val2",
	}
	result := ToGetParams(params)

	if result != "?test1=val1&test2=val2" {
		t.Fatalf("Expected ?test1=val1&test2=val2 but got %s", result)
	}
}

func TestGetWeatherdata(t *testing.T) {
	result, err := getWeatherData()

	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	fmt.Println(result.Name)
}

func TestGetAllEvents(t *testing.T) {
	result, err := GetAllExternalEvents()

	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	fmt.Println(result[0])
}
