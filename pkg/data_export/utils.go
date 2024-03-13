package internal_data

import (
	"fmt"
	"log"
	"os"

	"github.com/ClubCedille/hackqc2024/pkg/help"
)

func ConvertHelpsToGeoJSON(helps []*help.Help) ([]byte, error) {
	//todo if we want to export geojson
	return nil, nil
}


func fileSizeInMB(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file: %s", err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Printf("Error getting file info: %s", err)
	}

	size := fi.Size()
	
	if size < 1024 * 1024 {
		return "1" // lower than one Mo. API returns 409 error for values lower than 1.
	}
	return fmt.Sprintf("%.2f", float64(size)/1024/1024)
}