# Test Suite for HNG Step 0 API

This directory contains comprehensive tests for the HNG Step 0 API, including unit tests for helper functions and integration tests for API endpoints.

## Test Structure

### Files Overview

- `helpers_test.go` - Unit tests for helper functions
- `natural_language_test.go` - Tests for natural language filtering functionality
- `api_test.go` - Integration tests for API endpoints
- `test_helper.go` - Test utilities and setup functions

## Running Tests

### Run All Tests
```bash
go test ./tests/...
```

### Run Specific Test Files
```bash
# Test helper functions only
go test ./tests/helpers_test.go ./tests/test_helper.go

# Test natural language filtering only
go test ./tests/natural_language_test.go ./tests/test_helper.go

# Test API endpoints only
go test ./tests/api_test.go ./tests/test_helper.go
```

### Run Tests with Verbose Output
```bash
go test -v ./tests/...
```

### Run Tests with Coverage
```bash
go test -cover ./tests/...
```

## Test Categories

### 1. Helper Function Tests (`helpers_test.go`)

Tests all utility functions in the helpers package:

- **IsPalindrome**: Tests palindrome detection with various inputs
- **CountWords**: Tests word counting functionality
- **CountUniqueCharacters**: Tests unique character counting
- **CalculateSHA256**: Tests SHA256 hash calculation
- **CalculateCharacterFrequency**: Tests character frequency mapping
- **FindElement**: Tests element searching in arrays
- **StringApiHandler**: Tests the main handler methods

### 2. Natural Language Filtering Tests (`natural_language_test.go`)

Tests the natural language query parsing and filtering:

- **ParseNaturalLanguageQuery**: Tests query parsing for various natural language inputs
- **HasConflictingFilters**: Tests conflict detection in parsed filters
- **ApplyFilters**: Tests filter application to data
- **NaturalLanguageFilter**: Tests the basic natural language filtering

### 3. API Endpoint Tests (`api_test.go`)

Integration tests for all API endpoints:

- **Health Check**: Tests `/health` endpoint
- **POST /strings**: Tests string creation and analysis
- **GET /strings**: Tests string retrieval with filtering
- **GET /strings/:string_value**: Tests individual string retrieval
- **DELETE /strings/:string_value**: Tests string deletion
- **GET /strings/filter-by-natural-language**: Tests natural language filtering
- **GET /**: Tests API documentation endpoint

## Test Data

The tests use a global `TestBank` variable to store test data, which is reset between test runs to ensure test isolation.

## Expected Test Results

All tests should pass with the following coverage:

- Helper functions: 100% coverage
- Natural language filtering: 95%+ coverage
- API endpoints: 90%+ coverage

## Test Examples

### Example Test Run Output
```
=== RUN   TestIsPalindrome
--- PASS: TestIsPalindrome (0.00s)
=== RUN   TestCountWords
--- PASS: TestCountWords (0.00s)
=== RUN   TestParseNaturalLanguageQuery
--- PASS: TestParseNaturalLanguageQuery (0.00s)
=== RUN   TestHealthEndpoint
--- PASS: TestHealthEndpoint (0.00s)
=== RUN   TestPostStringsEndpoint
--- PASS: TestPostStringsEndpoint (0.00s)
PASS
ok      hng/step0/tests    0.123s
```

## Adding New Tests

When adding new tests:

1. Follow the existing naming convention: `TestFunctionName`
2. Use table-driven tests for multiple test cases
3. Reset test data between tests using `ResetTestBank()`
4. Include both positive and negative test cases
5. Test edge cases and error conditions

## Dependencies

The tests require:
- Go 1.19+
- Gin framework
- Standard Go testing package

## Continuous Integration

These tests are designed to run in CI/CD pipelines and should pass consistently across different environments.
