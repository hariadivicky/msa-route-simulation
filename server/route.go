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

// FindNext finds next coordinate point.
func (route Route) FindNext(currentPoint Point, reverse bool) (Point, error) {
	foundIndex := -1
	for i, point := range route.Points {
		if point.Latitude == currentPoint.Latitude && point.Longitude == currentPoint.Longitude {
			foundIndex = i
		}
	}

	if (foundIndex == len(route.Points)-1 && !reverse) || (foundIndex == 0 && reverse) {
		return Point{}, fmt.Errorf("point reach maximum value")
	}

	if foundIndex == -1 {
		return Point{}, fmt.Errorf("could not find next point")
	}

	if reverse {
		return route.Points[foundIndex-1], nil
	}

	return route.Points[foundIndex+1], nil
}

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
