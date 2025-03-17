package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mrjkey/chirpy/internal/database"
)

func main() {
	fmt.Println("Starting Server...")
	godotenv.Load()

	apicfg := apiConfig{}
	apicfg.fileserverHits.Store(0)
	apicfg.platform = os.Getenv("PLATFORM")

	dbURL := os.Getenv("DB_URL")
	fmt.Println(dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("unable to open db")
		os.Exit(1)
	}
	apicfg.db = database.New(db)
	// fmt.Println(dbQueries)

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	dir := http.Dir(".")

	fileserverHandler := http.StripPrefix("/app", http.FileServer(dir))
	mux.Handle("/app/", apicfg.middlewareMetricsInc(fileserverHandler))
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	// mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
	mux.HandleFunc("GET /admin/metrics", apicfg.handleMetrics())
	mux.HandleFunc("POST /admin/reset", apicfg.handleReset())
	mux.HandleFunc("POST /api/users", middlewareAddCfg(handleAddUser, &apicfg))

	mux.HandleFunc("GET /api/chirps", middlewareAddCfg(handleGetChirps, &apicfg))
	mux.HandleFunc("POST /api/chirps", middlewareAddCfg(handleAddChirp, &apicfg))

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

// func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
// 	type Chirp struct {
// 		Body string `json:"body"`
// 	}

// 	// type ChirpValid struct {
// 	// 	Valid bool `json:"valid"`
// 	// }
// 	type ReturnValue struct {
// 		CleanedBody string `json:"cleaned_body"`
// 	}
// 	decoder := json.NewDecoder(r.Body)
// 	chirp := Chirp{}
// 	err := decoder.Decode(&chirp)
// 	if err != nil {
// 		data := makeChirpError("could not decode incoming json")
// 		makeJsonResponse(w, data, http.StatusInternalServerError)
// 		return
// 	}

// 	if len(chirp.Body) > 140 {
// 		data := makeChirpError("Chirp is too long")
// 		makeJsonResponse(w, data, http.StatusBadRequest)
// 		return
// 	}

// 	cleanedBody := getCleanedBody(chirp.Body)

// 	returnValue := ReturnValue{CleanedBody: cleanedBody}
// 	data, _ := json.Marshal(returnValue)
// 	makeJsonResponse(w, data, http.StatusOK)
// }

func validateChirp(body string) (string, error) {

	if len(body) > 140 {
		return "", errors.New("chirp body is too long")
	}

	cleanedBody := getCleanedBody(body)

	return cleanedBody, nil
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

func middlewareAddCfg(
	function func(http.ResponseWriter, *http.Request, *apiConfig),
	cfg *apiConfig) func(http.ResponseWriter, *http.Request) {
	newFunction := func(w http.ResponseWriter, r *http.Request) {
		function(w, r, cfg)
	}
	return newFunction
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func handleAddUser(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	type NewUser struct {
		Email string `json:"email"`
	}

	user := NewUser{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		// do something
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), user.Email)
	if err != nil {
		fmt.Println(err)
		data := makeChirpError("could not create user in database")
		makeJsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(User(dbUser))
	if err != nil {
		data := makeChirpError("cannot marshel database user")
		makeJsonResponse(w, data, http.StatusInternalServerError)
		return
	}

	makeJsonResponse(w, data, http.StatusCreated)
}

func quickChirpError(w http.ResponseWriter, message string) {
	data := makeChirpError(message)
	makeJsonResponse(w, data, http.StatusInternalServerError)
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func handleAddChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		quickChirpError(w, err.Error())
	}
	body, err := validateChirp(chirp.Body)
	if err != nil {
		quickChirpError(w, err.Error())
	}

	args := database.AddChirpParams{
		Body:   body,
		UserID: chirp.UserID,
	}

	dbChirp, err := cfg.db.AddChirp(r.Context(), args)
	if err != nil {
		quickChirpError(w, err.Error())
	}
	data, err := json.Marshal(Chirp(dbChirp))
	if err != nil {
		quickChirpError(w, err.Error())
	}
	makeJsonResponse(w, data, http.StatusCreated)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		quickChirpError(w, err.Error())
	}

	respChirps := []Chirp{}
	for _, chirp := range chirps {
		respChirps = append(respChirps, Chirp(chirp))
	}

	data, err := json.Marshal(respChirps)
	if err != nil {
		quickChirpError(w, err.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(data)

}
