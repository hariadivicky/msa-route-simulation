package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// writeResponse is helper to marshal return data into http response body
func writeResponse(w http.ResponseWriter, data interface{}, statusCode ...int) {
	status := 200
	if len(statusCode) > 0 && statusCode[0] != 0 {
		status = statusCode[0]
	}

	b, err := json.Marshal(data)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(b)
}

func serveAPI(router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/routes", allowSimpleCors(getRoutes))
	subRouter.HandleFunc("/routes/{route}/coordinates", allowSimpleCors(getRouteCoordinates))
	subRouter.HandleFunc("/routes/{route}/start", allowSimpleCors(getRouteStart))
	subRouter.HandleFunc("/routes/{route}/next", allowSimpleCors(getNextCoordinate))
}

// allowSimpleCors allows cors request, only for reading data.
// this is helpful for development mode when our frontend running on different port.
func allowSimpleCors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	}

}

func getRoutes(w http.ResponseWriter, r *http.Request) {
	var routeNames []string
	for name := range routeMap {
		routeNames = append(routeNames, name)
	}

	writeResponse(w, routeNames)
}

func getRouteCoordinates(w http.ResponseWriter, r *http.Request) {
	routeName := mux.Vars(r)["route"]

	route, ok := routeMap[routeName]
	if !ok {
		writeResponse(w, "Route not found", http.StatusNotFound)
		return
	}

	writeResponse(w, route.Points)
}

func getRouteStart(w http.ResponseWriter, r *http.Request) {
	routeName := mux.Vars(r)["route"]
	reverse := r.URL.Query().Get("direction") == "reverse"
	route, ok := routeMap[routeName]
	if !ok {
		httpError(w, fmt.Errorf("route %s is not registered", routeName))
		return
	}

	var cursor int
	if reverse {
		// puts cursor to last point.
		cursor = len(route.Points) - 1
	}

	point := route.Points[cursor]
	writeResponse(w, point)

	// TODO: get route starting point, direction forward or reverse
}

func getNextCoordinate(w http.ResponseWriter, r *http.Request) {
	routeName := mux.Vars(r)["route"]
	reverse := r.URL.Query().Get("direction") == "reverse"
	currentPoint := r.URL.Query().Get("current")

	if currentPoint == "" {
		httpError(w, fmt.Errorf("current parameter is required"))
		return
	}

	route, ok := routeMap[routeName]
	if !ok {
		httpError(w, fmt.Errorf("route %s is not registered", routeName))
		return
	}

	var point Point

	if err := json.NewDecoder(strings.NewReader(currentPoint)).Decode(&point); err != nil {
		httpError(w, err)
		return
	}

	nextPoint, err := route.FindNext(point, reverse)
	if err != nil {
		httpError(w, err)
		return
	}

	writeResponse(w, nextPoint)
}
