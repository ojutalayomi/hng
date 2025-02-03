package main

import (
	"context"
	"encoding/json"
	"hng/step0/utils"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/rs/cors"
)

type HTTPResponseForIntro struct {
	Email            string `json:"email"`
	Current_datetime string `json:"current_datetime"`
	Github_url       string `json:"github_url"`
}

type HTTPResponseForClassifyNumber struct {
	Number      string   `json:"number"`
	Is_Prime    bool     `json:"is_prime"`
	Is_Perfect  bool     `json:"is_perfect"`
	Properties  []string `json:"properties"`
	Digital_Sum int      `json:"digital_sum"`
	Fun_Fact    string   `json:"fun_fact"`
}

type HTTPErrorResp struct {
	Number string `json:"number"`
	Error  bool   `json:"error"`
}

var (
	is_prime    bool
	is_perfect  bool
	properties  []string
	digital_sum int
	fun_fact    string
)

func Intro(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		http.Error(w, "missing option parameter", http.StatusBadRequest)
		return
	}

	now := time.Now()

	iso8601 := now.Format(time.RFC3339)

	log.Println(iso8601)

	json.NewEncoder(w).Encode(HTTPResponseForIntro{
		Email:            "ojutalayoayomide21@gmail.com",
		Current_datetime: iso8601,
		Github_url:       "https://github.com/ojutalayomi/hng",
	})
}

func ClassifyNumber(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	if r.Method != "GET" {
		http.Error(w, "missing option parameter", http.StatusBadRequest)
		return
	}

	initialNumber := r.URL.Query().Get("number")
	number, err := strconv.Atoi(initialNumber)
	if err != nil {
		errorResp := HTTPErrorResp{
			Number: initialNumber,
			Error:  true,
		}

		errorRespJSON, _ := json.Marshal(errorResp)
		http.Error(w, string(errorRespJSON), http.StatusBadRequest)
		return
	}

	// Check if number is prime
	is_prime = utils.IsPrime(number)

	// Check if number is perfect
	is_perfect = utils.IsPerfect(number)

	//Properties of the number
	properties = []string{}

	if utils.IsArmstrong(number) {
		properties = append(properties, "armstrong")
	}

	// Check if number is even or odd
	if utils.IsEven(number) {
		properties = append(properties, "even")
	} else {
		properties = append(properties, "odd")
	}

	digital_sum = utils.DigitalSum(number)

	fun_fact, _ = utils.FetchAPI("http://numbersapi.com/" + initialNumber + "/year?default=Boring+number+is+boring")

	json.NewEncoder(w).Encode(HTTPResponseForClassifyNumber{
		Number:      initialNumber,
		Is_Prime:    is_prime,
		Is_Perfect:  is_perfect,
		Properties:  properties,
		Digital_Sum: digital_sum,
		Fun_Fact:    fun_fact,
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

	mux.Handle("/api/classify-number", corsHandler(http.HandlerFunc(ClassifyNumber)))
	mux.Handle("/", corsHandler(http.HandlerFunc(Intro)))

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
