package internal_data

import (
	"fmt"
	"os"

	"github.com/ClubCedille/hackqc2024/pkg/help"
)

func ConvertHelpsToGeoJSON(helps []*help.Help) ([]byte, error) {
	//todo if we want to export geojson
	return nil, nil
}


func fileSizeInMB(filename string) string {
	fileInfo, _ := os.Stat(filename)
	fileSizeMB := fmt.Sprintf("%.0f", float64(fileInfo.Size()) / (1024.0 * 1024.0))
	return fileSizeMB
}