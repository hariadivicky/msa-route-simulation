package main

import (
	"io"
)

type csvReader struct{}

func (r *csvReader) Read(reader io.Reader) ([]Point, error) {
	// TODO: implement route csv reader
	return nil, nil
}

func createCSVReader() Reader {
	return &csvReader{}
}

func init() {
	RegisterReader("csv", createCSVReader)
}
