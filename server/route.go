package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Route contains route information
type Route struct {
	Name   string
	Points []Point
}

// Point represent route point
type Point struct {
	Longitude float64 `json:"lng"`
	Latitude  float64 `json:"lat"`
}

// Reader defines functionality route reader / parser
type Reader interface {
	Read(io.Reader) ([]Point, error)
}

// CreateReader is factory to create route reader
type CreateReader func() Reader

var (
	readers    map[string]CreateReader
	readerLock sync.Once
)

// RegisterReader add new route reader factory to list
func RegisterReader(name string, factory CreateReader) {
	readerLock.Do(func() {
		readers = make(map[string]CreateReader)
	})

	readers[name] = factory
}

// GetReader create reader base on given factory name
func GetReader(name string) Reader {
	factory, ok := readers[name]
	if !ok {
		fmt.Printf("%s is registered readers\n", name)
		return nil
	}

	return factory()
}

func parseRouteFile(file string, reader Reader) (Route, error) {
	f, err := os.Open(file)
	if err != nil {
		return Route{}, err
	}
	defer f.Close()

	coordinates, err := reader.Read(f)
	if err != nil {
		return Route{}, err
	}

	return Route{
		Name:   filepath.Base(file),
		Points: coordinates,
	}, nil
}

// parseRoutes lookup routes files on given root path
// read and parse into Routes
func parseRoutes(root string) []Route {
	var routes []Route

	filepath.Walk(root, func(p string, i os.FileInfo, err error) error {
		if err == nil && !i.IsDir() {
			// get file extension and check whether
			// route file is supported
			log.Println("found:", filepath.Base(p))
			reader := GetReader(strings.TrimLeft(filepath.Ext(p), "."))
			if reader != nil {
				route, err := parseRouteFile(p, reader)
				if err == nil {
					routes = append(routes, route)
				}
			}
		}

		return nil
	})

	return routes
}
