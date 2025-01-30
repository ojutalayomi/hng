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
  "github_url": "https://github.com/ojutalayomi/hng/tree/main/"
}
```

### Example Usage
```bash
curl http://localhost:8000/
```

## Related Resources
- [Hire Golang Developers](https://hng.tech/hire/golang-developers)