package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	ginMode    = os.Getenv("GIN_MODE")
	port       = os.Getenv("PORT")
	factApiUrl = os.Getenv("FACT_API_URL")
	userEmail  = os.Getenv("USER_EMAIL")
	userName   = os.Getenv("USER_NAME")
	userStack  = os.Getenv("USER_STACK")
)

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Stack string `json:"stack"`
}

type FactResponse struct {
	Fact string `json:"fact"`
}

type HTTPResponse struct {
	Status    string `json:"status"`
	User      User   `json:"user"`
	Timestamp string `json:"timestamp"`
	Fact      string `json:"fact"`
}

func FetchAPI(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Ensure response body is closed

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func Get(c *gin.Context) {

	c.Header("Content-Type", "application/json")
	if c.Request.Method != "GET" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing option parameter",
		})
		return
	}

	now := time.Now()

	iso8601 := now.Format(time.RFC3339)

	fact, err := FetchAPI(factApiUrl)
	if err != nil {
		log.Println(err)
	}

	var factResponse FactResponse
	err = json.Unmarshal(fact, &factResponse)
	if err != nil {
		log.Println(err)
	}

	if factResponse.Fact == "" {
		factResponse.Fact = "Unable to fetch fact"
	}

	c.JSON(http.StatusOK, HTTPResponse{
		Status: "success",
		User: User{
			Email: userEmail,
			Name:  userName,
			Stack: userStack,
		},
		Timestamp: iso8601,
		Fact:      factResponse.Fact,
	})
}

// setupRoutes configures all the API routes
func setupRoutes() *gin.Engine {
	// Create Gin router with default middleware (logger and recovery)
	router := gin.Default()
	fmt.Printf("\nGIN_MODE is %s\n", ginMode)
	gin.SetMode(ginMode)

	// Add CORS middleware to allow cross-origin requests
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "HNG Step 0 API",
			"timestamp": time.Now().UTC(),
		})
	})

	// Main endpoint for fetching link previews
	router.GET("/me", Get)

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		docs := map[string]any{
			"service":     "HNG Step 0 API",
			"version":     "1.0.0",
			"description": "API for fetching website metadata and me",
			"endpoints": map[string]any{
				"GET /me": map[string]any{
					"description": "Fetch me endpoint",
					"response": map[string]string{
						"email":     "Email",
						"name":      "Name",
						"status":    "Status",
						"fact":      "Fact",
						"timestamp": "Timestamp",
					},
				},
			},
		}

		c.JSON(http.StatusOK, docs)
	})

	return router
}

func main() {

	// Setup routes
	router := setupRoutes()

	// Get port from environment or use default
	port := fmt.Sprintf(":%s", port)

	fmt.Printf("🚀 HNG Step 0 API server starting on port %s\n", port)
	fmt.Println("📝 API Documentation available at: /")
	fmt.Println("🏥 Health check available at: /health")
	fmt.Println("🔗 Me endpoint: GET /me")

	// Start server
	if err := router.Run(port); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
	}
}
