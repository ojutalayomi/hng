# HNG Step 0 Project

## Description
This project is part of the HNG internship program, demonstrating basic API development skills using Go (Golang). It provides a simple endpoint that returns specific user information in JSON format.

## Setup Instructions

### Prerequisites
- Go (version 1.16 or later)
- Git

### Local Setup
1. Clone the repository:
```bash
git clone https://github.com/ojutalayomi/hng
cd step0
```

2. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8000`

## API Documentation

### Endpoint
- URL: `/`
- Method: `GET`

### Response Format
```json
{
  "email": "ojutalayoayomide21@gmail.com",
  "current_datetime": "2025-01-30T21:21:31+01:00",
  "github_url": "https://github.com/ojutalayomi/hng"
}
```

### Example Usage
```bash
curl http://localhost:8000/
```

### Task Description
An API that takes a number and returns interesting mathematical properties about it, along with a fun fact.

### Endpoint
- URL: `/api/classify-number?number={number}`
- Method: `GET`

### Response Format
```json
{
  "number": "371",
  "is_prime": false,
  "is_perfect": false,
  "properties": [
    "armstrong",
    "odd"
  ],
  "digital_sum": 11,
  "fun_fact": "371 is the year that Baekje forces storm the Goguryeo capital in P'yongyang (Korea)."
}
```

### Example Usage
```bash
curl http://localhost:8000/api/classify-number?number=371
```

## Related Resources
- [Hire Golang Developers](https://hng.tech/hire/golang-developers)