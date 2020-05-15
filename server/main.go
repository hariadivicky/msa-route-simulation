package main

import (
	"flag"
	"log"
)

var (
	routeMap map[string]Route
)

func main() {
	var (
		datadir  string
		assetdir string
		address  string
	)

	flag.StringVar(&datadir, "data-dir", "../data", "routes data directory")
	flag.StringVar(&assetdir, "asset-dir", "../ui/dist", "frontend assets directory")
	flag.StringVar(&address, "address", ":8000", "http address")

	flag.Parse()

	// parse coordinates files
	routeMap = make(map[string]Route)
	routes := parseRoutes(datadir)
	for _, route := range routes {
		routeMap[route.Name] = route
	}

	log.Fatal(startServer(ServerOption{
		Address:  address,
		AssetDir: assetdir,

		Routes: routes,
	}))
}
