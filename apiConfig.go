package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)


type apiConfig struct {
	fileserverHits atomic.Int32
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

func (cfg *apiConfig) handleMetricsReset() func(w http.ResponseWriter, r *http.Request) {
	function := func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset\n"))
	}
	return function
}
