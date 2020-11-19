package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonReader struct{}

func (r *jsonReader) Read(reader io.Reader) ([]Point, error) {
	var points []Point
	if err := json.NewDecoder(reader).Decode(&points); err != nil {
		return nil, fmt.Errorf("csvReader read error:  %v", err)
	}

	return points, nil
}

func createJSONReader() Reader {
	return &jsonReader{}
}

func init() {
	RegisterReader("json", createJSONReader)
}
