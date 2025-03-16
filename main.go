package main

import (
	"fmt"
	"net/http"
	"os"
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
		response := fmt.Sprintf("Hits: %v", hits)
		body := []byte(response)
		w.Header().Set("Content-Type", "text/plain")
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
		w.Write([]byte("Reset"))
	}
	return function
}

func main() {
	fmt.Println("Starting Server...")

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	dir := http.Dir(".")
	apicfg := apiConfig{}
	apicfg.fileserverHits.Store(0)

	fileserverHandler := http.StripPrefix("/app", http.FileServer(dir))
	mux.Handle("/app/", apicfg.middlewareMetricsInc(fileserverHandler))
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("GET /api/metrics", apicfg.handleMetrics())
	mux.HandleFunc("POST /api/reset", apicfg.handleMetricsReset())

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error starting server")
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := []byte("OK")
	w.Write(body)
}
