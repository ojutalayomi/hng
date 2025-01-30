package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"
)

type HTTPResponse struct {
	Email            string `json:"email"`
	Current_datetime string `json:"current_datetime"`
	Github_url       string `json:"github_url"`
}

func Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		http.Error(w, "missing option parameter", http.StatusBadRequest)
		return
	}

	now := time.Now()

	iso8601 := now.Format(time.RFC3339)

	log.Println(iso8601)

	json.NewEncoder(w).Encode(HTTPResponse{
		Email:            "ojutalayoayomide21@gmail.com",
		Current_datetime: iso8601,
		Github_url:       "https://github.com/ojutalayomi/hng/tree/main/",
	})
}

func main() {

	// Your existing server code
	mux := http.NewServeMux()

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	}).Handler

	mux.Handle("/api", corsHandler(http.HandlerFunc(Get)))
	mux.Handle("/", corsHandler(http.HandlerFunc(Get)))

	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down server...")

		// Create shutdown context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v\n", err)
		}
	}()

	log.Println("Serving at localhost:8080...")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
