package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/Sirupsen/logrus"
)

//Route struct describing a router route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
}

//Logger provides a decorator for logging
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)
		logrus.Infoln(
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

//NewRouter handles creation of new mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
