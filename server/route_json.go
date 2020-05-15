package main

import (
	"io"
)

type jsonReader struct{}

func (r *jsonReader) Read(reader io.Reader) ([]Point, error) {
	return nil, nil
}

func createJSONReader() Reader {
	return &jsonReader{}
}

func init() {
	RegisterReader("json", createJSONReader)
}
