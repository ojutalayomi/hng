package tests

import (
	helpers "hng/step0/helpers"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"racecar", true},
		{"Racecar", true},
		{"race car", true},
		{"A man a plan a canal Panama", true},
		{"hello", false},
		{"world", false},
		{"", true},  // empty string is considered palindrome
		{"a", true}, // single character is palindrome
		{"ab", false},
		{"aba", true},
		{"abba", true},
	}

	for _, test := range tests {
		result := helpers.IsPalindrome(test.input)
		if result != test.expected {
			t.Errorf("IsPalindrome(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello world", 2},
		{"hello", 1},
		{"", 0},
		{"hello world test", 3},
		{"  hello   world  ", 2}, // multiple spaces
		{"hello\nworld", 2},      // newline
		{"hello\tworld", 2},      // tab
	}

	for _, test := range tests {
		result := helpers.CountWords(test.input)
		if result != test.expected {
			t.Errorf("CountWords(%q) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestCountUniqueCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"hello", 4}, // h, e, l, o
		{"world", 5}, // w, o, r, l, d
		{"", 0},
		{"a", 1},
		{"aa", 1}, // only unique characters
		{"abcdef", 6},
		{"aabbcc", 3}, // a, b, c
	}

	for _, test := range tests {
		result := helpers.CountUniqueCharacters(test.input)
		if result != test.expected {
			t.Errorf("CountUniqueCharacters(%q) = %d, expected %d", test.input, result, test.expected)
		}
	}
}

func TestCalculateSHA256(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"a", "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
	}

	for _, test := range tests {
		result := helpers.CalculateSHA256(test.input)
		if result != test.expected {
			t.Errorf("CalculateSHA256(%q) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestCalculateCharacterFrequency(t *testing.T) {
	tests := []struct {
		input    string
		expected map[string]int
	}{
		{"hello", map[string]int{"h": 1, "e": 1, "l": 2, "o": 1}},
		{"", map[string]int{}},
		{"a", map[string]int{"a": 1}},
		{"aa", map[string]int{"a": 2}},
		{"aabbcc", map[string]int{"a": 2, "b": 2, "c": 2}},
	}

	for _, test := range tests {
		result := helpers.CalculateCharacterFrequency(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("CalculateCharacterFrequency(%q) length = %d, expected %d", test.input, len(result), len(test.expected))
		}
		for char, count := range test.expected {
			if result[char] != count {
				t.Errorf("CalculateCharacterFrequency(%q)[%s] = %d, expected %d", test.input, char, result[char], count)
			}
		}
	}
}

func TestFindElement(t *testing.T) {
	data := []helpers.Response{
		{Value: "hello", ID: "1"},
		{Value: "world", ID: "2"},
		{Value: "test", ID: "3"},
	}

	tests := []struct {
		targetValue string
		expected    int
	}{
		{"hello", 0},
		{"world", 1},
		{"test", 2},
		{"notfound", -1},
		{"", -1},
	}

	for _, test := range tests {
		result := helpers.FindElement(data, "value", test.targetValue)
		if result != test.expected {
			t.Errorf("FindElement(%q) = %d, expected %d", test.targetValue, result, test.expected)
		}
	}
}

func TestStringApiHandlerAnalyze(t *testing.T) {
	handler := helpers.StringApiHandler{String: "hello"}
	properties := handler.Analyze()

	if properties.Length != 5 {
		t.Errorf("Expected length 5, got %d", properties.Length)
	}
	if properties.IsPalindrome {
		t.Errorf("Expected IsPalindrome false, got true")
	}
	if properties.UniqueCharacters != 4 {
		t.Errorf("Expected UniqueCharacters 4, got %d", properties.UniqueCharacters)
	}
	if properties.WordCount != 1 {
		t.Errorf("Expected WordCount 1, got %d", properties.WordCount)
	}
	if properties.Sha256Hash == "" {
		t.Errorf("Expected non-empty SHA256Hash")
	}
}

func TestStringApiHandlerGetString(t *testing.T) {
	handler := helpers.StringApiHandler{String: "hello"}
	response := handler.GetString()

	if response.Value != "hello" {
		t.Errorf("Expected Value 'hello', got %s", response.Value)
	}
	if response.ID == "" {
		t.Errorf("Expected non-empty ID")
	}
	if response.CreatedAt == "" {
		t.Errorf("Expected non-empty CreatedAt")
	}
	if response.Properties.Length != 5 {
		t.Errorf("Expected Properties.Length 5, got %d", response.Properties.Length)
	}
}
