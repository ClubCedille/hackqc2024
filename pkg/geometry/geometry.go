package geometry

import (
	"strconv"
	"strings"

	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	pip "github.com/JamesLMilner/pip-go"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/xy"
)

func IsInGeom(point []float64, geometry mapobject.Geometry) bool {
	if point == nil || len(point) != 2 {
		return false
	}
	if geometry.GeomType == "Point" {
		const distCutoff = 0.01

		coord1 := geometry.AsGeomCoord()
		coord2 := geom.Coord(point[0:2])

		distance := xy.Distance(coord1, coord2)

		return distance < distCutoff
	} else if geometry.GeomType == "Polygon" {
		pipPoint := pip.Point{X: point[0], Y: point[1]}
		polygon := *geometry.AsPipPolygon()

		return pip.PointInPolygon(pipPoint, polygon)
	}

	return false
}

func ParseCoordinatesString(coordinates string) ([]float64, error) {
	coordinatesArray := strings.Split(coordinates, ",")

	var coordinatesArrayFloat []float64
	for i := len(coordinatesArray) - 1; i >= 0; i-- {
		coords, err := strconv.ParseFloat(strings.TrimSpace(coordinatesArray[i]), 64)
		if err != nil {
			return nil, err
		}
		coordinatesArrayFloat = append(coordinatesArrayFloat, coords)
	}

	return coordinatesArrayFloat, nil
}
