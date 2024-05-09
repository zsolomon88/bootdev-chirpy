package main

import "net/http"

func readinessHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
