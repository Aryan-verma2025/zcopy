package main

import (
	"fmt"
	"net/http"
)

func workingOnIt(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	fmt.Fprintf(w, "Working on it. Will be available soon")
}
