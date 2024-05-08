package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func main() {
	apiCfg := apiConfig{
		fileserverHits: 0,
	}
	httpMux := http.NewServeMux()
	httpMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	httpMux.HandleFunc("GET /healthz", readinessHandle)
	httpMux.HandleFunc("GET /metrics", apiCfg.metricsHandle)
	httpMux.HandleFunc("/reset", apiCfg.resetHandle)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	log.Println("Starting server on port: 8080")
	log.Fatal(httpServer.ListenAndServe())
}

func readinessHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (c *apiConfig) metricsHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hits: %d", c.fileserverHits)))
}

func (c *apiConfig) resetHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	c.fileserverHits = 0
}
