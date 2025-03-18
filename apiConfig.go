package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/mrjkey/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	tokenSecret    string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	function := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
	return function
}

func (cfg *apiConfig) handleMetrics() func(w http.ResponseWriter, r *http.Request) {
	function := func(w http.ResponseWriter, r *http.Request) {
		hits := cfg.fileserverHits.Load()

		response := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)
		body := []byte(response)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	}
	return function
}

func (cfg *apiConfig) handleReset() func(w http.ResponseWriter, r *http.Request) {
	function := func(w http.ResponseWriter, r *http.Request) {
		if cfg.platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		cfg.fileserverHits.Store(0)
		cfg.db.RemoveAllUsers(r.Context())
		cfg.db.RemoveChirps(r.Context())
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset\n"))
	}
	return function
}
