# Go-Vault

A simple and secure file storage web application built with Go, featuring user authentication and file management capabilities.

## Features

- User authentication with JWT tokens
- Secure file upload and storage
- File preview functionality
- SQLite database for user and file metadata
- Clean web interface built with Templ and Tailwind CSS
- Docker support for easy deployment

## Tech Stack

- **Backend**: Go with Echo framework
- **Database**: SQLite3
- **Frontend**: Templ templates, HTMX with Tailwind CSS
- **Authentication**: JWT tokens
- **Containerization**: Docker & Docker Compose

## Quick Start

### Using Docker Compose

1. Clone the repository
2. Set your JWT secret in environment variables:

   ```bash
   export JWT_SECRET="your-secure-secret-key"
   ```

3. Run with Docker Compose:

   ```bash
   docker-compose up -d
   ```

4. Access the application at `http://localhost:8080`

### Local Development

1. Install Go 1.25.5 or later
2. Install dependencies:

   ```bash
   go mod download
   ```

3. Run the application:

   ```bash
   go run cmd/main.go
   ```

## Project Structure

```
├── cmd/                 # Application entry point
├── internal/
│   ├── adapters/       # External adapters (handlers, repositories)
│   └── core/           # Business logic and services
├── templates/          # Templ templates for UI
├── static/             # Static assets
├── migrations/         # Database migrations
└── uploads/            # File storage directory
```

## Author

ELghaz
