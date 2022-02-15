package server

import (
	"compress/gzip"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"regexp"
	"strings"

	"github.com/axllent/adguard-home-bg/app"
	"github.com/axllent/adguard-home-bg/parser"
	"github.com/gorilla/mux"
)

var (
	// Listen tells the server on which port to listen on
	Listen = "0.0.0.0:8080" //string
)

//go:embed static
var embeddedFS embed.FS

// Start starts the HTTP server
func Start() error {

	listenRegex := regexp.MustCompile(`^(\b25[0-5]|\b2[0-4][0-9]|\b[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}:\d+$`)

	if !listenRegex.MatchString(Listen) {
		return fmt.Errorf("Invalid listen address: %s\n\nListen address should be in the format of <ip>:<port>", Listen)
	}

	serverRoot, err := fs.Sub(embeddedFS, "static")
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	r.HandleFunc("/", gzipHandlerFunc(parser.Build)).Methods("GET").Queries("url", "{filter}")
	r.PathPrefix("/").Handler(gzipHandler(http.FileServer(http.FS(serverRoot)))).Methods("GET")
	http.Handle("/", r)

	app.Log().InfoF("Starting server on http://%s", Listen)
	return http.ListenAndServe(Listen, nil)
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GzipHandlerFunc http middleware
func gzipHandlerFunc(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func gzipHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		h.ServeHTTP(gzipResponseWriter{Writer: gz, ResponseWriter: w}, r)
	})
}
