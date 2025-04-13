# Quiz App

A real-time quiz application with WebSocket support for interactive quizzes and leaderboards for test purposes.

## Features

- Real-time quiz participation via WebSockets
- Live leaderboard updates
- Multiple choice questions with time limits
- User score tracking
- Concurrent user support

## Architecture

The application follows a clean architecture pattern:

- **Domain Layer**: Core business entities and interfaces
- **Service Layer**: Business logic implementation
- **Repository Layer**: Data persistence
- **Handler Layer**: HTTP and WebSocket communication

## Tech Stack

- **Backend**: Go (Gin framework)
- **Frontend**: HTML, CSS, JavaScript
- **Communication**: WebSockets (gorilla/websocket)
- **Containerization**: Docker

## Project Structure

```
quiz-app/
├── cmd/
│   └── server/
├── internal/
│   ├── constants/
│   ├── domain/
│   │   ├── entities.go
│   │   ├── dtos.go
│   │   └── interfaces.go
│   ├── handler/         
│   │   ├── http_handler.go    # REST API handlers
│   │   └── websocket_handler.go # WebSocket communication
│   ├── repository/      
│   │   └── quiz_repository.go 
│   └── service/        
│       └── quiz_service.go  
└── static/              # Frontend assets
    ├── css/             
    ├── js/             
    └── index.html       
```

## Getting Started

### Prerequisites

- Go 1.23.3 or higher
- Docker (optional, for containerized deployment)

### Running Locally

1. Clone the repository:
   ```
   git clone https://github.com/nguyentanluctego/quiz-app
   cd quiz-app
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Run the application:
   ```
   go run cmd/server/main.go
   ```

4. Open your browser and navigate to `http://localhost:8080`

### Using Docker

1. Build the Docker image:
   ```
   docker build -t quiz-app .
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 quiz-app
   ```

3. Access the application at `http://localhost:8080`

## API Endpoints

### HTTP Endpoints

- `GET /api/quiz/:id` - Get quiz details

### WebSocket Endpoint

- `GET /ws` - WebSocket connection endpoint