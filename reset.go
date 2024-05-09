package main

import "net/http"

func (c *apiConfig) resetHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	c.fileserverHits = 0
}
