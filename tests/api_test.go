package tests

import (
	"bytes"
	"encoding/json"
	helpers "hng/step0/helpers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	router := SetupTestRouter()

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}
}

func TestPostStringsEndpoint(t *testing.T) {
	ResetTestBank()
	router := SetupTestRouter()

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectError    bool
	}{
		{
			"valid string",
			map[string]interface{}{
				"value": "hello world",
			},
			http.StatusOK,
			false,
		},
		{
			"empty string",
			map[string]interface{}{
				"value": "",
			},
			http.StatusUnprocessableEntity,
			true,
		},
		{
			"missing value field",
			map[string]interface{}{
				"other": "test",
			},
			http.StatusUnprocessableEntity,
			true,
		},
		{
			"invalid data type",
			map[string]interface{}{
				"value": 123,
			},
			http.StatusUnprocessableEntity,
			true,
		},
		{
			"duplicate string",
			map[string]interface{}{
				"value": "hello world",
			},
			http.StatusConflict,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(test.requestBody)
			req, _ := http.NewRequest("POST", "/strings", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, w.Code)
			}

			if test.expectError {
				var errorResponse map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
					t.Errorf("Failed to unmarshal error response: %v", err)
				}
				if errorResponse["error"] == nil {
					t.Errorf("Expected error in response")
				}
			} else {
				var response helpers.Response
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if response.Value != test.requestBody["value"] {
					t.Errorf("Expected value %v, got %v", test.requestBody["value"], response.Value)
				}
			}
		})
	}
}

func TestGetStringsEndpoint(t *testing.T) {
	ResetTestBank()
	router := SetupTestRouter()

	// First, add some test data
	testData := []string{"racecar", "hello world", "a", "abba", "long string here"}
	for _, value := range testData {
		jsonBody, _ := json.Marshal(map[string]string{"value": value})
		req, _ := http.NewRequest("POST", "/strings", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedCount  int
	}{
		{
			"get all strings",
			"",
			http.StatusOK,
			5,
		},
		{
			"filter by palindrome",
			"?is_palindrome=true",
			http.StatusOK,
			3, // racecar, a, abba
		},
		{
			"filter by word count",
			"?word_count=1",
			http.StatusOK,
			3, // racecar, a, abba
		},
		{
			"filter by min length",
			"?min_length=5",
			http.StatusOK,
			3, // racecar, hello world, long string here
		},
		{
			"filter by max length",
			"?max_length=4",
			http.StatusOK,
			2, // a, abba
		},
		{
			"filter by character",
			"?contains_character=a",
			http.StatusOK,
			3, // racecar, a, abba
		},
		{
			"invalid palindrome value",
			"?is_palindrome=invalid",
			http.StatusBadRequest,
			0,
		},
		{
			"invalid min_length",
			"?min_length=invalid",
			http.StatusBadRequest,
			0,
		},
		{
			"invalid word_count",
			"?word_count=invalid",
			http.StatusBadRequest,
			0,
		},
		{
			"invalid contains_character",
			"?contains_character=ab",
			http.StatusBadRequest,
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/strings"+test.queryParams, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, w.Code)
			}

			if test.expectedStatus == http.StatusOK {
				var response struct {
					Data  []helpers.Response `json:"data"`
					Count int                `json:"count"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if len(response.Data) != test.expectedCount {
					t.Errorf("Expected %d results, got %d", test.expectedCount, len(response.Data))
				}
			}
		})
	}
}

func TestGetStringByValueEndpoint(t *testing.T) {
	ResetTestBank()
	router := SetupTestRouter()

	// Add test data
	jsonBody, _ := json.Marshal(map[string]string{"value": "test string"})
	req, _ := http.NewRequest("POST", "/strings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	tests := []struct {
		name           string
		stringValue    string
		expectedStatus int
	}{
		{
			"existing string",
			"test string",
			http.StatusOK,
		},
		{
			"non-existing string",
			"non-existing",
			http.StatusNotFound,
		},
		{
			"empty string",
			"nonexistent",
			http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/strings/"+test.stringValue, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteStringEndpoint(t *testing.T) {
	ResetTestBank()
	router := SetupTestRouter()

	// Add test data
	jsonBody, _ := json.Marshal(map[string]string{"value": "delete test"})
	req, _ := http.NewRequest("POST", "/strings", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	tests := []struct {
		name           string
		stringValue    string
		expectedStatus int
	}{
		{
			"existing string",
			"delete test",
			http.StatusNoContent,
		},
		{
			"non-existing string",
			"non-existing",
			http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("DELETE", "/strings/"+test.stringValue, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, w.Code)
			}
		})
	}
}

func TestNaturalLanguageFilterEndpoint(t *testing.T) {
	ResetTestBank()
	router := SetupTestRouter()

	// Add test data
	testData := []string{"racecar", "hello world", "a", "abba", "long string here"}
	for _, value := range testData {
		jsonBody, _ := json.Marshal(map[string]string{"value": value})
		req, _ := http.NewRequest("POST", "/strings", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}

	tests := []struct {
		name           string
		query          string
		expectedStatus int
		expectedCount  int
	}{
		{
			"single word palindromic strings",
			"all single word palindromic strings",
			http.StatusOK,
			3, // racecar, a, abba
		},
		{
			"strings longer than 5 characters",
			"strings longer than 5 characters",
			http.StatusOK,
			3, // racecar, hello world, long string here
		},
		{
			"strings containing letter a",
			"strings containing the letter a",
			http.StatusOK,
			3, // racecar, a, abba
		},
		{
			"empty query",
			"",
			http.StatusBadRequest,
			0,
		},
		{
			"unparseable query",
			"some random unparseable text",
			http.StatusBadRequest,
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/strings/filter-by-natural-language?query="+test.query, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != test.expectedStatus {
				t.Errorf("Expected status %d, got %d", test.expectedStatus, w.Code)
			}

			if test.expectedStatus == http.StatusOK {
				var response struct {
					Data  []helpers.Response `json:"data"`
					Count int                `json:"count"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}
				if len(response.Data) != test.expectedCount {
					t.Errorf("Expected %d results, got %d", test.expectedCount, len(response.Data))
				}
			}
		})
	}
}

func TestAPIDocumentationEndpoint(t *testing.T) {
	router := SetupTestRouter()

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["service"] == nil {
		t.Errorf("Expected service field in response")
	}

	if response["endpoints"] == nil {
		t.Errorf("Expected endpoints field in response")
	}
}
