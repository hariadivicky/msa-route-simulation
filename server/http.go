package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

type ServerOption struct {
	Address  string
	AssetDir string

	Routes []Route
}

// helper to write error message and code to the response
func httpError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError

	if os.IsNotExist(err) {
		status = http.StatusNotFound
	}
	if os.IsPermission(err) {
		status = http.StatusForbidden
	}

	http.Error(w, http.StatusText(status), status)
}

// serve asset file on public path /assets
func serveAssetFile(router *mux.Router, assetDir string) {
	// server asset files
	fs := NoCache(http.FileServer(http.Dir(assetDir)))
	router.PathPrefix("/css").Handler(fs)
	router.PathPrefix("/fonts").Handler(fs)
	router.PathPrefix("/js").Handler(fs)

	// serve root frontend app (index.html)
	router.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		indexFile := filepath.Join(assetDir, "index.html")
		f, err := os.Open(indexFile)
		if err != nil {
			httpError(w, err)
			return
		}
		defer f.Close()

		d, err := f.Stat()
		if err != nil {
			httpError(w, err)
			return
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			httpError(w, err)
			return
		}

		w.Header().Set("Last-Modified", d.ModTime().UTC().Format(http.TimeFormat))
		if _, err := w.Write(b); err != nil {
			httpError(w, err)
		}
	})
}

func startServer(option ServerOption) error {
	// create http router
	router := mux.NewRouter()

	serveAPI(router)
	serveAssetFile(router, option.AssetDir)

	server := http.Server{
		Handler: router,
		Addr:    option.Address,
	}

	log.Println("Starting HTTP Server")
	log.Printf("HTTP Server is started on %s\n", option.Address)

	return server.ListenAndServe()
}

// NoCache wraps http handler to prevent client browser cache a request
// This middleware should be used only when using hot-reload option in local development
func NoCache(h http.Handler) http.Handler {
	var epoch = time.Unix(0, 0).Format(time.RFC1123)

	var noCacheHeaders = map[string]string{
		"Expires":         epoch,
		"Cache-Control":   "no-cache, private, max-age=0",
		"Pragma":          "no-cache",
		"X-Accel-Expires": "0",
	}

	var etagHeaders = []string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}

	fn := func(w http.ResponseWriter, r *http.Request) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if r.Header.Get(v) != "" {
				r.Header.Del(v)
			}
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
