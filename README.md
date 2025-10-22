# HNG Step 0 API

This project provides a string analysis API with endpoints to analyze strings, filter data using natural language or structured queries, and more.

## Features

- Add and manage analyzed strings
- Query string properties: palindrome, word count, length, contains character, etc.
- Filter strings using both URL query params and natural language syntax
- REST API with JSON input/output
- Written in Go, using Gin web framework

## API Endpoints

### Health Check

```
GET /health
```

Returns API status and timestamp.

### Add String

```
POST /strings
Content-Type: application/json

{
  "value": "your string here"
}
```

Adds a string and analyzes its properties. Returns 409 Conflict if already present.

### List & Filter Strings

```
GET /strings
```

Optional query params (structured):

- `is_palindrome=true|false`
- `word_count=N`
- `min_length=N`
- `max_length=N`
- `contains_character=a`

Example:

```
GET /strings?is_palindrome=true&min_length=3
```

Returns all matching strings and their properties.

### Natural Language Filter

```
GET /strings/filter-by-natural-language?query=palindromes longer than 4 letters
```

Uses natural phrases to interpret filter queries. Returns:

- Filtered data,
- Count,
- How the query was interpreted.

### API Documentation

```
GET /
```

Returns summary of available endpoints.

## Running the API

1. Install Go and [Gin](https://github.com/gin-gonic/gin)
2. Clone the repo
3. Run:

```
go run main.go
```

The service listens on the port set by the `PORT` environment variable (default: 8080), and obeys the `GIN_MODE` environment variable for running in debug or release.

## Running the Tests

Run all tests:

```
go test ./tests/...
```

## Example Filters

- **Structured:** `GET /strings?word_count=2`
- **Natural Language:** `GET /strings/filter-by-natural-language?query=palindromes only`

## Example Response (POST /strings)

```json
{
  "value": "racecar",
  "properties": {
    "is_palindrome": true,
    "word_count": 1,
    "length": 7
  }
}
```

---
