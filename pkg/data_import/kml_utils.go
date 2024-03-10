package data_import

import (
	"archive/zip"
	"encoding/xml"
	"os"
)

func ExtractKMLFile(src string, dst string) error {
	//Unzip KMZ file
	archive, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		if file.Name != "doc.kml" {
			continue
		}
		open, err := file.Open()
		if err != nil {
			return err
		}
		defer open.Close()
		create, err := os.Create(dst)
		if err != nil {
			return err
		}
		defer create.Close()
		create.ReadFrom(open)
	}
	return nil
}

func LoadKmlFile(fileName string) (*KML, error) {
	kmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer kmlFile.Close()

	kml := KML{}
	err = xml.NewDecoder(kmlFile).Decode(&kml)
	if err != nil {
		return nil, err
	}

	return &kml, nil
}

type KML struct {
	XMLName  xml.Name    `xml:"kml"`
	XMLNS    string      `xml:"xmlns,attr"`
	GX       string      `xml:"gx,attr"`
	Document KMLDocument `xml:"Document"`
}

type KMLDocument struct {
	XMLName   xml.Name       `xml:"Document"`
	Placemark []KMLPlacemark `xml:"Placemark"`
}

type KMLPlacemark struct {
	Centroid string `xml:"ExtendData>Data>value"`
	Polygon  string `xml:"Polygon"`
}
