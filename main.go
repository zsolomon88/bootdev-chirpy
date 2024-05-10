package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits int
	jwtSecret      string
}

func main() {

	godotenv.Load()

	apiCfg := apiConfig{
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	httpMux.HandleFunc("GET /api/healthz", readinessHandle)
	httpMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandle)
	httpMux.HandleFunc("GET /api/reset", apiCfg.resetHandle)
	httpMux.HandleFunc("POST /api/chirps", createHandle)
	httpMux.HandleFunc("GET /api/chirps/{chirpId}", getHandle)
	httpMux.HandleFunc("GET /api/chirps", getHandle)
	httpMux.HandleFunc("POST /api/users", apiCfg.createUserHandle)
	httpMux.HandleFunc("POST /api/login", apiCfg.authenticateHandle)
	httpMux.HandleFunc("PUT /api/users", apiCfg.updateUsrHandle)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	log.Println("Starting server on port: 8080")
	log.Fatal(httpServer.ListenAndServe())
}
