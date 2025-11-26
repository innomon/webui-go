# Go Backend Implementation Document

## 1. Overview

This document provides a technical overview of the Go backend that was migrated from the original Python implementation. The goal of this migration was to create a more performant, statically-typed, and maintainable backend service while preserving the core functionalities of the original application.

The backend is a standalone Go application that serves a RESTful API and handles real-time communication via Socket.IO.

## 2. Technology Stack

- **Language:** Go
- **Web Framework:** Chi (v5)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Real-time Communication:** Socket.IO for Go
- **Authentication:** JWT (JSON Web Tokens)

## 3. Project Structure

The `backend/` directory is organized into several packages to separate concerns and improve code organization.

```
backend/
├───config/
│   └───config.go        # Environment variable management
├───database/
│   └───database.go      # Database connection and schema migration
├───handlers/
│   ├───auth.go          # User registration and login handlers
│   ├───chat.go          # Chat and message handlers
│   ├───file.go          # File and folder management handlers
│   ├───knowledge.go     # Knowledge base handlers
│   ├───llm.go           # LLM interaction handlers
│   ├───model.go         # Model management handlers
│   ├───prompt.go        # Prompt management handlers
│   ├───tool.go          # Tool management handlers
│   └───user_admin.go    # User administration handlers
├───middleware/
│   └───auth.go          # JWT authentication middleware
├───models/
│   ├───chat.go          # Chat and Message data models
│   ├───file.go          # File and Folder data models
│   ├───knowledge.go     # Knowledge Base data models
│   ├───llm.go           # Ollama and OpenAI request/response structs
│   ├───model.go         # AI Model data models
│   ├───prompt.go        # Prompt data models
│   ├───tool.go          # Tool data models
│   ├───user.go          # User data model
│   └───user_admin.go    # Structs for user administration forms
├───routes/
│   ├───chat.go          # Chat API routes definition
│   ├───file.go          # File and Folder API routes
│   ├───knowledge.go     # Knowledge Base API routes
│   ├───llm.go           # LLM API routes
│   ├───model.go         # Model management API routes
│   ├───prompt.go        # Prompt management API routes
│   ├───tool.go          # Tool management API routes
│   └───user_admin.go    # User administration API routes
├───services/
│   └───llm.go           # Services for calling Ollama and OpenAI APIs
├───utils/
│   └───response.go      # Utility functions for API responses
├───.env                   # Local environment variables (DB connection, etc.)
├───go.mod                 # Go module dependencies
├───go.sum                 # Go module checksums
└───main.go                # Application entry point
```

## 4. Getting Started

### 4.1. Database Setup

The backend requires a PostgreSQL database. A local instance can be easily run using Docker:

```bash
# Create a persistent volume for the data
docker volume create postgres_data

# Run the PostgreSQL container
docker run --name postgres-db \
  -e POSTGRES_USER=admin \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=webui \
  -p 5432:5432 \
  -v postgres_data:/var/lib/postgresql/data \
  -d postgres
```

### 4.2. Environment Variables

Create a `.env` file inside the `backend/` directory with the following content:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=password
DB_NAME=webui
DB_SSLMODE=disable

# Example LLM configurations
OLLAMA_BASE_URL=http://localhost:11434
OPENAI_API_BASE_URL=https://api.openai.com
OPENAI_API_KEY=your_openai_api_key
```

### 4.3. Running the Server

Navigate to the `backend/` directory and run the main application:

```bash
go run main.go
```

The server will start on `http://localhost:8080`. Upon startup, it will automatically connect to the database and run all necessary schema migrations.

## 5. Key Features Implemented

- **Full User Authentication:** Registration, login, and protected routes using JWT.
- **Complete Chat API:** CRUD for chats and messages.
- **Real-time Chat:** Socket.IO integration for broadcasting new messages to participants in a chat room.
- **LLM Integration:** Handlers and services to connect to both Ollama and OpenAI compatible APIs.
- **File and Folder Management:** API for uploading, downloading, and organizing files.
- **Knowledge Base Management:** Basic CRUD for creating and managing knowledge bases and associating files with them.
- **Model Management:** CRUD for managing AI model configurations.
- **Prompt Management:** CRUD for creating, retrieving, and managing reusable prompts.
- **Tool Management:** Basic CRUD for managing external tools.
- **User Administration:** Basic endpoints for listing, updating, and deleting users.

