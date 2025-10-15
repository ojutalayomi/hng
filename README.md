# HNG Step 0 Project

## Description
This project is part of the HNG internship program, demonstrating basic API development skills using Go (Golang). It provides a simple endpoint that returns specific user information in JSON format, as well as a random cat fact.

## Setup Instructions

### Prerequisites
- Go (version 1.16 or later)
- Git

### Environment Variables
- `GIN_MODE`: Set to `release` for production.
- `PORT`: Set to the port you want to run the server on.
- `FACT_API_URL`: Set to the URL of the fact API you want to use.
- `USER_EMAIL`: Set to the email of the user you want to return.
- `USER_NAME`: Set to the name of the user you want to return.
- `USER_STACK`: Set to the stack of the user you want to return.

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

The server will start on `http://localhost:8080`

## API Documentation

### Endpoints

#### `GET /me`
Returns your profile information with a cat fact.

**Example Response**
```json
{
  "status": "success",
  "user": {
    "email": "user@gmail.com",
    "name": "User Name",
    "stack": "User Stack"
  },
  "timestamp": "2025-01-30T21:21:31+01:00",
  "fact": "A group of cats is called a clowder."
}
```

#### `GET /health`
Returns server health status.

**Example Response**
```json
{
  "status": "healthy",
  "service": "HNG Step 0 API",
  "timestamp": "2025-01-30T21:21:31.123Z"
}
```

## Technologies Used

- [Go (Golang)](https://golang.org/)
- [Gin Web Framework](https://github.com/gin-gonic/gin)

## Project Structure

- `main.go`: Entry point, server setup, routes, and controllers.

## Notes

- All CORS headers are set to allow requests from any origin.
- Timestamps are provided in ISO 8601 format (RFC3339).
- The `/me` endpoint fetches a fact from [catfact.ninja](https://catfact.ninja).
- To change the port, edit the `port` variable in `main.go`.

## Example Usage

```bash
curl http://localhost:8080/me
```

## License

MIT

