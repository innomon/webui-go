# Backend Specification Document

## 1. Overview

This document outlines the backend API for the WebUI-Go application. It has been updated to reflect the successful migration from Python to Go. The backend provides a RESTful API and a Socket.IO interface for real-time communication.

## 2. Technology Stack

- **Language:** Go
- **Web Framework:** Chi (v5)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT (JSON Web Tokens)
- **Real-time Communication:** Socket.IO

## 3. API Endpoints

The following table lists the final API endpoints for the Go backend.

| Category      | Endpoint                         | HTTP Method | Description                                       |
|---------------|----------------------------------|-------------|---------------------------------------------------|
| **Auths**     | `/api/auth/login`                | `POST`      | Authenticate a user and receive a JWT.            |
|               | `/api/auth/register`             | `POST`      | Register a new user.                              |
| **Chats**     | `/api/chats`                     | `GET`       | Get a list of chats for the authenticated user.   |
|               | `/api/chats`                     | `POST`      | Create a new chat.                                |
|               | `/api/chats/{id}/messages`       | `GET`       | Get all messages for a specific chat.             |
|               | `/api/chats/{id}/messages`       | `POST`      | Send a message in a chat and trigger LLM response.|
| **LLMs**      | `/api/chat/completions`          | `POST`      | Get a completion from an LLM (Ollama/OpenAI).     |
| **Files**     | `/api/files/upload`              | `POST`      | Upload a file.                                    |
|               | `/api/files/{id}`                | `GET`       | Get file metadata by ID.                          |
|               | `/api/files/{id}/download`       | `GET`       | Download a file's content.                        |
|               | `/api/files/{id}`                | `DELETE`    | Delete a file.                                    |
| **Folders**   | `/api/folders`                   | `GET`       | Get top-level folders or folders by `parent_id`.  |
|               | `/api/folders`                   | `POST`      | Create a new folder.                              |
|               | `/api/folders/{id}`              | `GET`       | Get a folder's content (files and subfolders).    |
|               | `/api/folders/{id}`              | `DELETE`    | Delete a folder.                                  |
| **Knowledge** | `/api/knowledge/create`          | `POST`      | Create a new knowledge base.                      |
|               | `/api/knowledge`                 | `GET`       | Get all knowledge bases for the user.             |
|               | `/api/knowledge/{id}`            | `GET`       | Get a specific knowledge base by ID.              |
|               | `/api/knowledge/{id}`            | `PUT`       | Update a knowledge base.                          |
|               | `/api/knowledge/{id}`            | `DELETE`    | Delete a knowledge base.                          |
|               | `/api/knowledge/{id}/file/add`   | `POST`      | Add a file to a knowledge base.                   |
|               | `/api/knowledge/{id}/file/remove`| `POST`      | Remove a file from a knowledge base.              |
| **Models**    | `/api/models/create`             | `POST`      | Create a new model configuration.                 |
|               | `/api/models/list`               | `GET`       | Get a list of available models.                   |
|               | `/api/models/{id}`               | `GET`       | Get a model by ID.                                |
|               | `/api/models/{id}`               | `PUT`       | Update a model.                                   |
|               | `/api/models/{id}`               | `DELETE`    | Delete a model.                                   |
| **Prompts**   | `/api/prompts/create`            | `POST`      | Create a new prompt.                              |
|               | `/api/prompts`                   | `GET`       | Get all prompts for the user.                     |
|               | `/api/prompts/command/{command}` | `GET`       | Get a prompt by its command (e.g., `/summarize`). |
|               | `/api/prompts/command/{command}/update` | `PUT`  | Update a prompt by its command.                  |
|               | `/api/prompts/command/{command}/delete` | `DELETE`| Delete a prompt by its command.                 |
| **Tools**     | `/api/tools/create`              | `POST`      | Create a new tool.                                |
|               | `/api/tools`                     | `GET`       | Get all tools for the user.                       |
|               | `/api/tools/id/{id}`             | `GET`       | Get a tool by ID.                                 |
|               | `/api/tools/id/{id}/update`      | `PUT`       | Update a tool.                                    |
|               | `/api/tools/id/{id}/delete`      | `DELETE`    | Delete a tool.                                    |
| **Users**     | `/api/users`                     | `GET`       | **[Admin]** Get a list of all users.              |
|               | `/api/users/{id}`                | `PUT`       | **[Admin]** Update a user's details.              |
|               | `/api/users/{id}`                | `DELETE`    | **[Admin]** Delete a user.                        |
|               | `/api/user/me`                   | `GET`       | Get the current authenticated user's profile.     |

## 4. Data Models (GORM)

- **User:** Stores user information, including `ID`, `Email`, `Password` (hashed), `Name`, and `Role`.
- **Chat:** Represents a conversation, with a `UserID` and `Title`.
- **Message:** A single message within a `Chat`, containing `Role` (e.g., "user", "assistant") and `Content`.
- **File:** Metadata for an uploaded file, including `Name`, `Path`, `MimeType`, and `Size`.
- **Folder:** A directory to organize files, with support for nested structures (`ParentID`).
- **Knowledge:** A knowledge base, with a `Name`, `Description`, and a JSONB array of `FileIDs`.
- **Model:** Configuration for an AI model, with `ID`, `Name`, `Meta` (JSONB), and `Params` (JSONB).
- **Prompt:** A reusable prompt with a `Title`, `Content`, and a unique `Command` (e.g., `/summarize`).
- **Tool:** A custom tool with `ID`, `Name`, `Content` (Python code), and `Specs` (JSONB).

## 5. Real-time Communication

Real-time communication is handled via Socket.IO. After establishing a connection, the client must authenticate.

- **Authentication:**
  - Client emits `auth` with the JWT token.
  - Server responds with `authenticated` on success or `authError` on failure.
- **Chat:**
  - Client emits `joinChat` with a `chatID` to join a room.
  - Server emits `message` to all clients in a room when a new message is created (both from the user and from the LLM).
  - Client can emit `leaveChat` to exit a room.

## 6. Getting Started

1. **Install Go and Docker.**
2. **Run a PostgreSQL container** (see `backend/IMPLEMENTATION.md` for details).
3. **Create a `backend/.env` file** (see `backend/IMPLEMENTATION.md` for details).
4. **Run the backend server:** `cd backend && go run main.go`
