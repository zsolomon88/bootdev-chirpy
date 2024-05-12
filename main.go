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
	polkaKey       string
}

func main() {

	godotenv.Load()

	apiCfg := apiConfig{
		fileserverHits: 0,
		jwtSecret:      os.Getenv("JWT_SECRET"),
		polkaKey:       os.Getenv("POLKA_KEY"),
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/app/*", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	httpMux.HandleFunc("GET /api/healthz", readinessHandle)
	httpMux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandle)
	httpMux.HandleFunc("GET /api/reset", apiCfg.resetHandle)
	httpMux.HandleFunc("POST /api/chirps", apiCfg.createHandle)
	httpMux.HandleFunc("GET /api/chirps/{chirpId}", getHandle)
	httpMux.HandleFunc("DELETE /api/chirps/{chirpId}", apiCfg.deleteHandle)
	httpMux.HandleFunc("GET /api/chirps", getHandle)
	httpMux.HandleFunc("POST /api/users", apiCfg.createUserHandle)
	httpMux.HandleFunc("POST /api/login", apiCfg.authenticateHandle)
	httpMux.HandleFunc("PUT /api/users", apiCfg.updateUsrHandle)
	httpMux.HandleFunc("POST /api/refresh", apiCfg.refreshHandle)
	httpMux.HandleFunc("POST /api/revoke", apiCfg.revokeTokenHandle)
	httpMux.HandleFunc("POST /api/polka/webhooks", apiCfg.redWebhook)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	log.Println("Starting server on port: 8080")
	log.Fatal(httpServer.ListenAndServe())
}
