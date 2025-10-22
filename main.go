package main

import (
	"encoding/json"
	"fmt"
	helpers "hng/step0/helpers"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	ginMode = os.Getenv("GIN_MODE")
	port    = os.Getenv("PORT")
	bank    []helpers.Response
)

// setupRoutes configures all the API routes
func SetupRoutes() *gin.Engine {
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
			"service":   "HNG Step 1 API",
			"timestamp": time.Now().UTC(),
		})
	})

	// Main endpoint for creating/analyzing strings
	router.POST("/strings", func(c *gin.Context) {
		var requestBody struct {
			Value string `json:"value" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				if ute.Value != "string" {
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": fmt.Sprintf("Invalid data type for \"%s\" (must be string)", ute.Field)})
					return
				}
			} else {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request body or missing \"value\" field"})
				return
			}
		}

		if helpers.FindElement(bank, "value", requestBody.Value) != -1 {
			c.JSON(http.StatusConflict, gin.H{"error": "String already exists in the system"})
			return
		}

		handler := helpers.StringApiHandler{String: requestBody.Value}
		response := handler.GetString()
		bank = append(bank, response)

		c.JSON(http.StatusCreated, response)
	})

	router.GET("/strings/:string_value", func(c *gin.Context) {
		stringValue := c.Param("string_value")
		index := helpers.FindElement(bank, "value", stringValue)
		if index == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusOK, bank[index])
	})

	router.GET("/strings", func(c *gin.Context) {
		isPalindrome := c.DefaultQuery("is_palindrome", "")
		minLength := c.DefaultQuery("min_length", "0")
		maxLength := c.DefaultQuery("max_length", "0")
		wordCount := c.DefaultQuery("word_count", "0")
		containsCharacter := c.DefaultQuery("contains_character", "")
		filteredBank := make([]helpers.Response, 0)

		fmt.Println(isPalindrome, minLength, maxLength, wordCount, containsCharacter)
		var filteredResponse struct {
			Data           []helpers.Response `json:"data"`
			Count          int                `json:"count"`
			FiltersApplied map[string]string  `json:"filters_applied"`
		}

		filteredResponse.FiltersApplied = make(map[string]string)

		// Validate query parameter values and types, respond 400 if invalid
		if isPalindrome != "" && isPalindrome != "true" && isPalindrome != "false" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter values or types"})
			return
		}
		if minLength != "0" && minLength != "" {
			if _, err := strconv.Atoi(minLength); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter values or types"})
				return
			}
		}
		if maxLength != "0" && maxLength != "" {
			if _, err := strconv.Atoi(maxLength); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter values or types"})
				return
			}
		}
		if wordCount != "0" && wordCount != "" {
			if _, err := strconv.Atoi(wordCount); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter values or types"})
				return
			}
		}
		if containsCharacter != "" {
			if len(containsCharacter) != 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter values or types"})
				return
			}
		}

		if isPalindrome == "" && minLength == "0" && maxLength == "0" && wordCount == "0" && containsCharacter == "" {
			c.JSON(http.StatusOK, filteredResponse)
			return
		}

		for _, item := range bank {
			match := true

			if isPalindrome == "true" && isPalindrome != "" {
				filteredResponse.FiltersApplied["is_palindrome"] = isPalindrome
				if !helpers.IsPalindrome(item.Value) {
					match = false
				}
			} else if isPalindrome == "false" && isPalindrome != "" {
				filteredResponse.FiltersApplied["is_palindrome"] = isPalindrome
				if helpers.IsPalindrome(item.Value) {
					match = false
				}
			}

			if minLength != "0" {
				filteredResponse.FiltersApplied["min_length"] = minLength
				minLengthInt, err := strconv.Atoi(minLength)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid min_length parameter"})
					return
				}
				if len(item.Value) < minLengthInt {
					match = false
				}
			}

			if maxLength != "0" {
				filteredResponse.FiltersApplied["max_length"] = maxLength
				maxLengthInt, err := strconv.Atoi(maxLength)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid max_length parameter"})
					return
				}
				if len(item.Value) > maxLengthInt {
					match = false
				}
			}

			if wordCount != "0" {
				filteredResponse.FiltersApplied["word_count"] = wordCount
				wordCountInt, err := strconv.Atoi(wordCount)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word_count parameter"})
					return
				}
				if helpers.CountWords(item.Value) != wordCountInt {
					match = false
				}
			}

			if containsCharacter != "" {
				filteredResponse.FiltersApplied["contains_character"] = containsCharacter
				if !strings.ContainsRune(item.Value, rune(containsCharacter[0])) {
					match = false
				}
			}

			if match {
				filteredBank = append(filteredBank, item)
			}
		}

		filteredResponse.Data = filteredBank
		filteredResponse.Count = len(filteredBank)

		c.JSON(http.StatusOK, filteredResponse)
	})

	router.GET("/strings/filter-by-natural-language", func(c *gin.Context) {
		query := c.DefaultQuery("query", "")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
			return
		}

		// Parse the natural language query
		parsedFilters, err := helpers.ParseNaturalLanguageQuery(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse natural language query"})
			return
		}

		// Check for conflicting filters
		if helpers.HasConflictingFilters(parsedFilters) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Query parsed but resulted in conflicting filters"})
			return
		}

		// Apply filters to get matching strings
		filteredBank := helpers.ApplyFilters(bank, parsedFilters)

		// Build response
		var filteredResponse struct {
			Data             []helpers.Response `json:"data"`
			Count            int                `json:"count"`
			InterpretedQuery struct {
				Original      string                 `json:"original"`
				ParsedFilters map[string]interface{} `json:"parsed_filters"`
			} `json:"interpreted_query"`
		}

		filteredResponse.Data = filteredBank
		filteredResponse.Count = len(filteredBank)
		filteredResponse.InterpretedQuery.Original = query
		filteredResponse.InterpretedQuery.ParsedFilters = parsedFilters

		c.JSON(http.StatusOK, filteredResponse)
	})

	router.DELETE("/strings/:string_value", func(c *gin.Context) {
		stringValue := c.Param("string_value")
		index := helpers.FindElement(bank, "value", stringValue)
		if index == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		bank = append(bank[:index], bank[index+1:]...)
		c.JSON(http.StatusNoContent, nil)
	})

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		docs := map[string]any{
			"service":     "HNG Step 0 API",
			"version":     "1.0.0",
			"description": "API for string analysis and natural language filtering",
			"endpoints": map[string]any{
				"POST /strings": map[string]any{
					"description": "Create/Analyze Strings endpoint",
					"request": map[string]string{
						"value": "string to analyze",
					},
					"response": helpers.Response{
						ID:    "sha256_hash_value",
						Value: "string to analyze",
						Properties: helpers.PropertiesMap{
							Length:           17,
							IsPalindrome:     false,
							UniqueCharacters: 12,
							WordCount:        3,
							Sha256Hash:       "abc123...",
							CharacterFrequencyMap: helpers.CharacterFrequencyMap{
								"s": 2,
								"t": 3,
								"r": 2,
							},
						},
						CreatedAt: "2025-08-27T10:00:00Z",
					},
				},
				"GET /strings": map[string]any{
					"description": "Get all strings with optional filtering",
					"query_params": map[string]string{
						"is_palindrome":      "true/false",
						"min_length":         "number",
						"max_length":         "number",
						"word_count":         "number",
						"contains_character": "single character",
					},
				},
				"GET /strings/filter-by-natural-language": map[string]any{
					"description": "Filter strings using natural language queries",
					"query_params": map[string]string{
						"query": "natural language query (e.g., 'all single word palindromic strings')",
					},
					"example_queries": []string{
						"all single word palindromic strings",
						"strings longer than 10 characters",
						"palindromic strings that contain the first vowel",
						"strings containing the letter z",
					},
				},
				"GET /strings/:string_value": map[string]any{
					"description": "Get a specific string by value",
				},
				"DELETE /strings/:string_value": map[string]any{
					"description": "Delete a specific string by value",
				},
			},
		}

		c.JSON(http.StatusOK, docs)
	})

	return router
}

func main() {
	router := SetupRoutes()

	// Get port from environment or use default
	port := fmt.Sprintf(":%s", port)

	fmt.Printf("🚀 HNG Step 1 API server starting on port %s\n", port)
	fmt.Println("📝 API Documentation available at: /")
	fmt.Println("🏥 Health check available at: /health")
	fmt.Println("🔗 Me endpoint: GET /me")

	// Start server
	if err := router.Run(port); err != nil {
		fmt.Printf("❌ Failed to start server: %v\n", err)
	}
}
