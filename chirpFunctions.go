package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mrjkey/chirpy/internal/auth"
	"github.com/mrjkey/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func convertChirp(dbChirp database.Chirp) Chirp {
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}
	return chirp
}

func handleAddChirp(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		quickChirpError(w, err.Error())
		return
	}
	userID, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		errData := makeChirpError(err.Error())
		makeJsonResponse(w, errData, http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err = decoder.Decode(&chirp)
	if err != nil {
		quickChirpError(w, err.Error())
		return
	}

	body, err := validateChirp(chirp.Body)
	if err != nil {
		quickChirpError(w, err.Error())
		return
	}

	args := database.AddChirpParams{
		Body:   body,
		UserID: userID,
	}

	dbChirp, err := cfg.db.AddChirp(r.Context(), args)
	if err != nil {
		quickChirpError(w, err.Error())
		return
	}
	data, err := json.Marshal(convertChirp(dbChirp))
	if err != nil {
		quickChirpError(w, err.Error())
		return
	}
	makeJsonResponse(w, data, http.StatusCreated)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		quickChirpError(w, err.Error())
		return
	}

	respChirps := []Chirp{}
	for _, chirp := range chirps {
		respChirps = append(respChirps, convertChirp(chirp))
	}

	data, err := json.Marshal(respChirps)
	if err != nil {
		quickChirpError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(data)

}

func handlGetChirpById(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	id := r.PathValue("chirpID")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		errData := makeChirpError(err.Error())
		makeJsonResponse(w, errData, http.StatusNotFound)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), parsedId)
	if err != nil {
		errData := makeChirpError(err.Error())
		makeJsonResponse(w, errData, http.StatusNotFound)
		return
	}

	data, err := json.Marshal(Chirp(chirp))
	if err != nil {
		errData := makeChirpError(err.Error())
		makeJsonResponse(w, errData, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
