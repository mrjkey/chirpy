package main

import (
	"encoding/json"
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
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apicfg.handleMetrics())
	mux.HandleFunc("POST /admin/reset", apicfg.handleMetricsReset())

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

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}

	type ChirpValid struct {
		Valid bool `json:"valid"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		dat := makeChirpError("could not decode incoming json")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(dat)
		return
	}

	if len(chirp.Body) > 140 {
		dat := makeChirpError("Chirp is too long")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
		return
	}

	valid := ChirpValid{Valid: true}
	dat, _ := json.Marshal(valid)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func makeChirpError(text string) []byte {
	type ChirpError struct {
		Error string `json:"error"`
	}
	c_err := ChirpError{Error: text}
	dat, err := json.Marshal(c_err)
	if err != nil {
		fmt.Println("unable to marshal chirp error")
		dat = []byte{}
	}
	return dat
}
