package tests

import (
	helpers "hng/step0/helpers"
	"testing"
)

func TestParseNaturalLanguageQuery(t *testing.T) {
	tests := []struct {
		query    string
		expected map[string]interface{}
		hasError bool
	}{
		{
			"all single word palindromic strings",
			map[string]interface{}{
				"word_count":    1,
				"is_palindrome": true,
			},
			false,
		},
		{
			"strings longer than 10 characters",
			map[string]interface{}{
				"min_length": 11,
			},
			false,
		},
		{
			"palindromic strings that contain the first vowel",
			map[string]interface{}{
				"is_palindrome":      true,
				"contains_character": "a",
			},
			false,
		},
		{
			"strings containing the letter z",
			map[string]interface{}{
				"contains_character": "z",
			},
			false,
		},
		{
			"two word strings",
			map[string]interface{}{
				"word_count": 2,
			},
			false,
		},
		{
			"strings shorter than 5 characters",
			map[string]interface{}{
				"max_length": 4,
			},
			false,
		},
		{
			"non-palindromic strings",
			map[string]interface{}{
				"is_palindrome": false,
			},
			false,
		},
		{
			"strings containing the letter a",
			map[string]interface{}{
				"contains_character": "a",
			},
			false,
		},
		{
			"unparseable query",
			nil,
			true,
		},
		{
			"",
			nil,
			true,
		},
	}

	for _, test := range tests {
		result, err := helpers.ParseNaturalLanguageQuery(test.query)

		if test.hasError {
			if err == nil {
				t.Errorf("ParseNaturalLanguageQuery(%q) expected error, got nil", test.query)
			}
			continue
		}

		if err != nil {
			t.Errorf("ParseNaturalLanguageQuery(%q) unexpected error: %v", test.query, err)
			continue
		}

		if len(result) != len(test.expected) {
			t.Errorf("ParseNaturalLanguageQuery(%q) length = %d, expected %d", test.query, len(result), len(test.expected))
		}

		for key, expectedValue := range test.expected {
			if result[key] != expectedValue {
				t.Errorf("ParseNaturalLanguageQuery(%q)[%s] = %v, expected %v", test.query, key, result[key], expectedValue)
			}
		}
	}
}

func TestHasConflictingFilters(t *testing.T) {
	tests := []struct {
		filters  map[string]interface{}
		expected bool
	}{
		{
			map[string]interface{}{
				"min_length": 10,
				"max_length": 5,
			},
			true,
		},
		{
			map[string]interface{}{
				"min_length": 5,
				"max_length": 10,
			},
			false,
		},
		{
			map[string]interface{}{
				"is_palindrome": true,
				"word_count":    1,
			},
			false,
		},
		{
			map[string]interface{}{
				"min_length": 10,
				"max_length": 10,
			},
			false,
		},
		{
			map[string]interface{}{
				"word_count": 1,
			},
			false,
		},
	}

	for _, test := range tests {
		result := helpers.HasConflictingFilters(test.filters)
		if result != test.expected {
			t.Errorf("HasConflictingFilters(%v) = %v, expected %v", test.filters, result, test.expected)
		}
	}
}

func TestApplyFilters(t *testing.T) {
	// Create test data
	data := []helpers.Response{
		{Value: "racecar", Properties: helpers.PropertiesMap{IsPalindrome: true, WordCount: 1}},
		{Value: "hello world", Properties: helpers.PropertiesMap{IsPalindrome: false, WordCount: 2}},
		{Value: "a", Properties: helpers.PropertiesMap{IsPalindrome: true, WordCount: 1}},
		{Value: "abba", Properties: helpers.PropertiesMap{IsPalindrome: true, WordCount: 1}},
		{Value: "long string here", Properties: helpers.PropertiesMap{IsPalindrome: false, WordCount: 3}},
	}

	tests := []struct {
		name     string
		filters  map[string]interface{}
		expected int // expected number of matches
	}{
		{
			"palindrome filter",
			map[string]interface{}{
				"is_palindrome": true,
			},
			3, // racecar, a, abba
		},
		{
			"word count filter",
			map[string]interface{}{
				"word_count": 1,
			},
			3, // racecar, a, abba
		},
		{
			"combined filters",
			map[string]interface{}{
				"is_palindrome": true,
				"word_count":    1,
			},
			3, // racecar, a, abba
		},
		{
			"min length filter",
			map[string]interface{}{
				"min_length": 5,
			},
			3, // racecar, hello world, long string here
		},
		{
			"max length filter",
			map[string]interface{}{
				"max_length": 4,
			},
			2, // a, abba
		},
		{
			"character filter",
			map[string]interface{}{
				"contains_character": "a",
			},
			3, // racecar, a, abba
		},
		{
			"no matches",
			map[string]interface{}{
				"contains_character": "z",
			},
			0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := helpers.ApplyFilters(data, test.filters)
			if len(result) != test.expected {
				t.Errorf("ApplyFilters(%s) returned %d results, expected %d", test.name, len(result), test.expected)
			}
		})
	}
}

func TestNaturalLanguageFilter(t *testing.T) {
	data := []helpers.Response{
		{Value: "hello world"},
		{Value: "test"},
		{Value: "another test"},
	}

	tests := []struct {
		query    string
		expected int
	}{
		{"hello", 1},
		{"test", 2},
		{"world", 1},
		{"nonexistent", 0},
		{"", 3}, // Empty string matches all strings
	}

	for _, test := range tests {
		result := helpers.NaturalLanguageFilter(data, test.query)
		if len(result) != test.expected {
			t.Errorf("NaturalLanguageFilter(%q) returned %d results, expected %d", test.query, len(result), test.expected)
		}
	}
}
