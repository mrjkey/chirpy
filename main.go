package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting Server...")

	mux := http.NewServeMux()
	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	dir := http.Dir(".")
	fileserverHandler := http.StripPrefix("/app", http.FileServer(dir))
	mux.Handle("/app/", fileserverHandler)
	mux.HandleFunc("/healthz", handleHealthz)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("error starting server")
		fmt.Println(err)
		os.Exit(1)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	body := []byte("OK")
	w.Write(body)
}
