# GEMINI.md

## Project Overview

This is a real-time chat application built with Go. It allows users to create and join chat rooms to communicate in real-time.

The application follows a standard Go web application architecture:
- **Backend:** Written in Go, using the standard `net/http` library for the web server and `golang.org/x/net/websocket` for real-time communication.
- **Frontend:** Server-side rendered using Go's `html/template` package.
- **Dependencies:** It uses `github.com/gorilla/sessions` for managing user sessions.
- **Structure:** The code is organized into several packages:
    - `cmd`: Contains the main application entry point.
    - `routes`: Defines the HTTP routes and handlers.
    - `server`: Manages the WebSocket server and connections.
    - `types`: Defines the data structures used throughout the application.

## Building and Running

### Standard Execution

To run the application, use the following command from the project root:

```sh
go run cmd/main.go
```

The server will start on `http://localhost:3000`.

### Development with Live Reload

This project is configured to use the `air` tool for live reloading during development. To use it, first make sure you have `air` installed:

```sh
go install github.com/air-verse/air@latest
```

Then, run the application with:

```sh
air
```

This will automatically watch for file changes and restart the server.

## Development Conventions

- **Dependency Management:** The project uses Go modules. To add a new dependency, use `go get`.
- **Routing:** Routes are defined in `routes/routes.go` and registered in the `RegisterRoutes` function.
- **WebSocket Handling:** The core WebSocket logic is in `server/server.go`. The `HandleWS` function manages new connections, and the `broadcast` function sends messages to clients in a specific room.
- **Data Types:** Shared data structures are defined in the `types` package.
