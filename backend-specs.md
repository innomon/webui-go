# Backend Specification Document

## 1. Overview

This document outlines the backend API for the WebUI-Go application. It is intended to provide a clear specification for the API endpoints that the frontend will consume. The backend is built with Go and provides a RESTful API for the frontend.

## 2. Technology Stack

- **Language:** Go
- **Framework:** Net/HTTP (standard library)
- **Database:** (To be determined)

## 3. API Endpoints

The following table lists the API endpoints that are expected to be used by the frontend, based on the structure of the `src/lib/apis` directory. The exact routes and HTTP methods are inferred and may need to be adjusted based on the frontend's implementation.

| Category      | Endpoint                         | HTTP Method | Description                                       |
|---------------|----------------------------------|-------------|---------------------------------------------------|
| **Audio**     | `/api/audio`                     | `POST`      | Upload an audio file.                             |
|               | `/api/audio/{id}`                | `GET`       | Get an audio file by ID.                          |
| **Auths**     | `/api/auth/login`                | `POST`      | Authenticate a user.                              |
|               | `/api/auth/logout`               | `POST`      | Log out a user.                                   |
|               | `/api/auth/register`             | `POST`      | Register a new user.                              |
| **Channels**  | `/api/channels`                  | `GET`       | Get a list of channels.                           |
|               | `/api/channels`                  | `POST`      | Create a new channel.                             |
|               | `/api/channels/{id}`             | `GET`       | Get a channel by ID.                              |
|               | `/api/channels/{id}`             | `PUT`       | Update a channel.                                 |
|               | `/api/channels/{id}`             | `DELETE`    | Delete a channel.                                 |
| **Chats**     | `/api/chats`                     | `GET`       | Get a list of chats.                              |
|               | `/api/chats`                     | `POST`      | Create a new chat.                                |
|               | `/api/chats/{id}`                | `GET`       | Get a chat by ID.                                 |
|               | `/api/chats/{id}/messages`       | `GET`       | Get messages for a chat.                          |
|               | `/api/chats/{id}/messages`       | `POST`      | Send a message in a chat.                         |
| **Configs**   | `/api/configs`                   | `GET`       | Get application configuration.                    |
|               | `/api/configs`                   | `PUT`       | Update application configuration.                 |
| **Files**     | `/api/files`                     | `POST`      | Upload a file.                                    |
|               | `/api/files/{id}`                | `GET`       | Get a file by ID.                                 |
| **Folders**   | `/api/folders`                   | `GET`       | Get a list of folders.                            |
|               | `/api/folders`                   | `POST`      | Create a new folder.                              |
| **Models**    | `/api/models`                    | `GET`       | Get a list of available models.                   |
| **Notes**     | `/api/notes`                     | `GET`       | Get a list of notes.                              |
|               | `/api/notes`                     | `POST`      | Create a new note.                                |
|               | `/api/notes/{id}`                | `GET`       | Get a note by ID.                                 |
|               | `/api/notes/{id}`                | `PUT`       | Update a note.                                    |
|               | `/api/notes/{id}`                | `DELETE`    | Delete a note.                                    |
| **Users**     | `/api/users`                     | `GET`       | Get a list of users.                              |
|               | `/api/users/{id}`                | `GET`       | Get a user by ID.                                 |

## 4. Real-time Communication

In addition to the REST API, the backend will also provide real-time communication using Socket.IO and Server-Sent Events (SSE) for features like chat and notifications. The exact implementation of these real-time features will be documented separately.

## 5. Getting Started

1. **Install Go:** [https://golang.org/doc/install](https://golang.org/doc/install)
2. **Run the backend server:** `go run backend/main.go`
