package tests

import (
	helpers "hng/step0/helpers"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// TestBank is a global variable to store test data
var TestBank []helpers.Response

// SetupTestRouter creates a test router with the same routes as main
func SetupTestRouter() *gin.Engine {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

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

	// POST /strings endpoint
	router.POST("/strings", func(c *gin.Context) {
		var requestBody struct {
			Value string `json:"value" binding:"required"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request body or missing \"value\" field"})
			return
		}

		if helpers.FindElement(TestBank, "value", requestBody.Value) != -1 {
			c.JSON(http.StatusConflict, gin.H{"error": "String already exists in the system"})
			return
		}

		handler := helpers.StringApiHandler{String: requestBody.Value}
		response := handler.GetString()
		TestBank = append(TestBank, response)

		c.JSON(http.StatusOK, response)
	})

	// GET /strings/:string_value endpoint
	router.GET("/strings/:string_value", func(c *gin.Context) {
		stringValue := c.Param("string_value")
		index := helpers.FindElement(TestBank, "value", stringValue)
		if index == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		c.JSON(http.StatusOK, TestBank[index])
	})

	// GET /strings endpoint with filtering
	router.GET("/strings", func(c *gin.Context) {
		isPalindrome := c.DefaultQuery("is_palindrome", "")
		minLength := c.DefaultQuery("min_length", "0")
		maxLength := c.DefaultQuery("max_length", "0")
		wordCount := c.DefaultQuery("word_count", "0")
		containsCharacter := c.DefaultQuery("contains_character", "")
		filteredBank := make([]helpers.Response, 0)

		var filteredResponse struct {
			Data           []helpers.Response `json:"data"`
			Count          int                `json:"count"`
			FiltersApplied map[string]string  `json:"filters_applied"`
		}

		filteredResponse.FiltersApplied = make(map[string]string)

		// Validate query parameters
		if isPalindrome != "" && isPalindrome != "true" && isPalindrome != "false" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter values or types"})
			return
		}

		for _, item := range TestBank {
			match := true

			if isPalindrome == "true" {
				filteredResponse.FiltersApplied["is_palindrome"] = isPalindrome
				if !helpers.IsPalindrome(item.Value) {
					match = false
				}
			} else if isPalindrome == "false" {
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
				if len(containsCharacter) != 1 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contains_character parameter"})
					return
				}
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

	// DELETE /strings/:string_value endpoint
	router.DELETE("/strings/:string_value", func(c *gin.Context) {
		stringValue := c.Param("string_value")
		index := helpers.FindElement(TestBank, "value", stringValue)
		if index == -1 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		TestBank = append(TestBank[:index], TestBank[index+1:]...)
		c.JSON(http.StatusNoContent, nil)
	})

	// Natural language filter endpoint
	router.GET("/strings/filter-by-natural-language", func(c *gin.Context) {
		query := c.DefaultQuery("query", "")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
			return
		}

		parsedFilters, err := helpers.ParseNaturalLanguageQuery(query)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse natural language query"})
			return
		}

		if helpers.HasConflictingFilters(parsedFilters) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Query parsed but resulted in conflicting filters"})
			return
		}

		filteredBank := helpers.ApplyFilters(TestBank, parsedFilters)

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

	// API documentation endpoint
	router.GET("/", func(c *gin.Context) {
		docs := map[string]any{
			"service":     "HNG Step 0 API",
			"version":     "1.0.0",
			"description": "API for string analysis and natural language filtering",
			"endpoints": map[string]any{
				"GET /health": "Health check endpoint",
			},
		}
		c.JSON(http.StatusOK, docs)
	})

	return router
}

// ResetTestBank clears the test bank
func ResetTestBank() {
	TestBank = []helpers.Response{}
}
