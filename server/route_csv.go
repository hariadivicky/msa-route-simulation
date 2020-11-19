package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type csvReader struct{}

// Read implements io.Reader, converts csv into Point slices.
func (r *csvReader) Read(reader io.Reader) ([]Point, error) {

	lines, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("csvReader error reading: %v", err)
	}

	var points []Point
	for line, cols := range lines {
		if len(cols) != 2 {
			return nil, fmt.Errorf("csvReader error reading: invalid column count on line %d", line)
		}

		lng, err := strconv.ParseFloat(cols[0], 64)
		if err != nil {
			return nil, fmt.Errorf("csvReader parsing longitude error: %v", err)
		}

		lat, err := strconv.ParseFloat(cols[1], 64)
		if err != nil {
			return nil, fmt.Errorf("csvReader parsing latitude error: %v", err)
		}

		points = append(points, Point{lng, lat})
	}

	return points, nil
}

func createCSVReader() Reader {
	return &csvReader{}
}

func init() {
	RegisterReader("csv", createCSVReader)
}
