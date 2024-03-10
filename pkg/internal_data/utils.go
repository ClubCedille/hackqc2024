package internal_data

import (
	"fmt"
	"os"

	"github.com/ClubCedille/hackqc2024/pkg/event"
)

func ConvertEventsToGeoJSON(events []*event.Event) ([]byte, error) {
	//todo once we know why and what we want to export specifically
	return nil, nil
}


func fileSizeInMB(filename string) string {
	fileInfo, _ := os.Stat(filename)
	fileSizeMB := fmt.Sprintf("%.0f", float64(fileInfo.Size()) / (1024.0 * 1024.0))
	return fileSizeMB
}