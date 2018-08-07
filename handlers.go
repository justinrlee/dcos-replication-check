package main

import (
	"encoding/json"
	// "io"
	// "io/ioutil"
	// "log"
	"net/http"
)

//Index for slash, returns version
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode("Health check tool running")
}
