package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StringApiHandler struct {
	String string
}

type CharacterFrequencyMap map[string]int

type PropertiesMap struct {
	Length                int                   `json:"length"`
	IsPalindrome          bool                  `json:"is_palindrome"`
	UniqueCharacters      int                   `json:"unique_characters"`
	WordCount             int                   `json:"word_Count"`
	Sha256Hash            string                `json:"sha256_hash"`
	CharacterFrequencyMap CharacterFrequencyMap `json:"character_frequency_map"`
}

type Response struct {
	ID         string        `json:"id"`
	Value      string        `json:"value"`
	Properties PropertiesMap `json:"properties"`
	CreatedAt  string        `json:"created_at"`
}

func (h *StringApiHandler) Analyze() PropertiesMap {
	return PropertiesMap{
		Length:                len(h.String),
		IsPalindrome:          IsPalindrome(h.String),
		UniqueCharacters:      CountUniqueCharacters(h.String),
		WordCount:             CountWords(h.String),
		Sha256Hash:            CalculateSHA256(h.String),
		CharacterFrequencyMap: CalculateCharacterFrequency(h.String),
	}
}

func (h *StringApiHandler) GetString() Response {
	properties := h.Analyze()
	return Response{
		ID:         properties.Sha256Hash,
		Value:      h.String,
		Properties: properties,
		CreatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
}

func (h *StringApiHandler) GetAllStrings(c *gin.Context) {}

func (h *StringApiHandler) NaturalLanguageFilter(c *gin.Context) {}

func (h *StringApiHandler) DeleteString(c *gin.Context) {}

func CountCharacters(s string) int {
	count := len(s)
	return count
}

func CountWords(s string) int {
	words := strings.Fields(s)
	return len(words)
}
func IsPalindrome(s string) bool {
	string := strings.ToLower(s)
	string = strings.ReplaceAll(string, " ", "") // Normalize the string
	length := len(string)
	for i := 0; i < length/2; i++ {
		if string[i] != string[length-i-1] {
			return false
		}
	}
	return true
}

func CountUniqueCharacters(s string) int {
	charSet := make(map[rune]struct{})
	for _, char := range s {
		charSet[char] = struct{}{}
	}
	return len(charSet)
}

func CalculateSHA256(s string) string {
	hash := sha256.Sum256([]byte(s))

	hashHex := hex.EncodeToString(hash[:])
	return hashHex
}

func CalculateCharacterFrequency(s string) map[string]int {
	frequency := make(map[string]int)
	for _, char := range s {
		frequency[string(char)]++
	}
	return frequency
}

func FindElement(data []Response, targetKey string, targetValue string) int {
	return slices.IndexFunc(data, func(r Response) bool {
		return r.Value == targetValue
	})
}

func NaturalLanguageFilter(data []Response, naturalLanguageFilter string) []Response {
	var filtered []Response
	for _, r := range data {
		if strings.Contains(r.Value, naturalLanguageFilter) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

// ParseNaturalLanguageQuery parses natural language queries into structured filters
func ParseNaturalLanguageQuery(query string) (map[string]interface{}, error) {
	query = strings.ToLower(query)
	filters := make(map[string]interface{})

	// Check for palindrome-related queries
	if strings.Contains(query, "palindromic") || strings.Contains(query, "palindrome") {
		if strings.Contains(query, "not") || strings.Contains(query, "non") {
			filters["is_palindrome"] = false
		} else {
			filters["is_palindrome"] = true
		}
	}

	// Check for word count queries
	if strings.Contains(query, "single word") || strings.Contains(query, "one word") {
		filters["word_count"] = 1
	} else if strings.Contains(query, "two word") || strings.Contains(query, "2 word") {
		filters["word_count"] = 2
	} else if strings.Contains(query, "three word") || strings.Contains(query, "3 word") {
		filters["word_count"] = 3
	}

	// Check for length-related queries
	if strings.Contains(query, "longer than") {
		// Extract number from "longer than X characters"
		parts := strings.Split(query, "longer than")
		if len(parts) > 1 {
			numberPart := strings.TrimSpace(parts[1])
			numberPart = strings.Split(numberPart, " ")[0] // Get just the number
			if num, err := strconv.Atoi(numberPart); err == nil {
				filters["min_length"] = num + 1
			}
		}
	} else if strings.Contains(query, "shorter than") {
		// Extract number from "shorter than X characters"
		parts := strings.Split(query, "shorter than")
		if len(parts) > 1 {
			numberPart := strings.TrimSpace(parts[1])
			numberPart = strings.Split(numberPart, " ")[0] // Get just the number
			if num, err := strconv.Atoi(numberPart); err == nil {
				filters["max_length"] = num - 1
			}
		}
	}

	// Check for character-specific queries
	if strings.Contains(query, "contain") {
		// Look for specific character mentions
		if strings.Contains(query, "letter z") || strings.Contains(query, "character z") {
			filters["contains_character"] = "z"
		} else if strings.Contains(query, "letter a") || strings.Contains(query, "character a") {
			filters["contains_character"] = "a"
		} else if strings.Contains(query, "letter e") || strings.Contains(query, "character e") {
			filters["contains_character"] = "e"
		} else if strings.Contains(query, "letter i") || strings.Contains(query, "character i") {
			filters["contains_character"] = "i"
		} else if strings.Contains(query, "letter o") || strings.Contains(query, "character o") {
			filters["contains_character"] = "o"
		} else if strings.Contains(query, "letter u") || strings.Contains(query, "character u") {
			filters["contains_character"] = "u"
		} else if strings.Contains(query, "first vowel") {
			filters["contains_character"] = "a"
		}
	}

	// If no filters were found, return an error
	if len(filters) == 0 {
		return nil, fmt.Errorf("unable to parse query")
	}

	return filters, nil
}

// HasConflictingFilters checks if the parsed filters have any conflicts
func HasConflictingFilters(filters map[string]interface{}) bool {
	// Check for conflicting palindrome settings
	if val, exists := filters["is_palindrome"]; exists {
		if val == false {
			// If explicitly set to false, check if we also have palindrome=true somewhere
			// This is a basic check - you might want to expand this logic
		}
	}

	// Check for conflicting length settings
	if minLen, minExists := filters["min_length"]; minExists {
		if maxLen, maxExists := filters["max_length"]; maxExists {
			if minLen.(int) > maxLen.(int) {
				return true
			}
		}
	}

	return false
}

// ApplyFilters applies the parsed filters to the data
func ApplyFilters(data []Response, filters map[string]interface{}) []Response {
	var filtered []Response

	for _, item := range data {
		match := true

		// Apply palindrome filter
		if isPalindrome, exists := filters["is_palindrome"]; exists {
			if isPalindrome.(bool) != IsPalindrome(item.Value) {
				match = false
			}
		}

		// Apply word count filter
		if wordCount, exists := filters["word_count"]; exists {
			if CountWords(item.Value) != wordCount.(int) {
				match = false
			}
		}

		// Apply minimum length filter
		if minLength, exists := filters["min_length"]; exists {
			if len(item.Value) < minLength.(int) {
				match = false
			}
		}

		// Apply maximum length filter
		if maxLength, exists := filters["max_length"]; exists {
			if len(item.Value) > maxLength.(int) {
				match = false
			}
		}

		// Apply character filter
		if containsChar, exists := filters["contains_character"]; exists {
			if !strings.ContainsRune(item.Value, rune(containsChar.(string)[0])) {
				match = false
			}
		}

		if match {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
