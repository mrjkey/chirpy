package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	_ "github.com/lib/pq"
	"github.com/mrjkey/chirpy/internal/database"
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

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	fmt.Println(dbQueries)

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

	err = server.ListenAndServe()
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

	// type ChirpValid struct {
	// 	Valid bool `json:"valid"`
	// }
	type ReturnValue struct {
		CleanedBody string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		data := makeChirpError("could not decode incoming json")
		makeJsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	if len(chirp.Body) > 140 {
		data := makeChirpError("Chirp is too long")
		makeJsonResponse(w, data, http.StatusBadRequest)
		return
	}

	cleanedBody := getCleanedBody(chirp.Body)

	returnValue := ReturnValue{CleanedBody: cleanedBody}
	data, _ := json.Marshal(returnValue)
	makeJsonResponse(w, data, http.StatusOK)
}

func getCleanedBody(body string) string {
	var profanity = map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := profanity[lowerWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func makeJsonResponse(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
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
