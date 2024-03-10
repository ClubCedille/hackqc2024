package data_import_test

import (
	"testing"

	"github.com/ClubCedille/hackqc2024/pkg/data_import"
)

func TestLoadKML(t *testing.T) {
	kml, err := data_import.LoadKmlFile(`C:\Users\Antoine\Documents\Programming\hackqc2024\tmp\outageAreas20240310183020.kml`)

	if err != nil {
		t.Error(err)
	}

	t.Log(kml)
}
