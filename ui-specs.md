# UI Specification Document

## 1. Overview

This document outlines the architecture and implementation details of the WebUI-Go frontend. It is intended to help developers understand the project structure, conventions, and key technologies used.

The application is a Single-Page Application (SPA) built with **SvelteKit**, a modern Svelte framework. It utilizes **Tailwind CSS** for styling and **i18next** for internationalization. A unique aspect of this project is its use of **Pyodide** and **ONNX Runtime**, enabling client-side Python execution and machine learning model inference directly in the browser.

## 2. Technology Stack

- **Framework:** SvelteKit
- **Language:** TypeScript
- **Styling:** Tailwind CSS
- **UI Components:** Svelte components
- **State Management:** Svelte Stores
- **Internationalization:** i18next
- **Client-side Python:** Pyodide
- **Client-side ML:** ONNX Runtime
- **Real-time Communication:** Socket.IO, Server-Sent Events
- **API Communication:** REST (via `fetch` API)

## 3. Project Structure

The project follows the standard SvelteKit directory structure:

```
/
├── src/
│   ├── lib/
│   │   ├── apis/         # Backend communication layer
│   │   ├── components/   # Reusable Svelte components
│   │   ├── i18n/         # i18next configuration and locales
│   │   ├── stores/       # Svelte stores for state management
│   │   ├── types/        # TypeScript type definitions
│   │   └── utils/        # Utility functions
│   ├── routes/         # File-based routing
│   │   ├── (app)/      # Authenticated routes
│   │   └── auth/       # Authentication routes
│   └── app.html        # Main HTML entry point
├── static/             # Static assets (images, fonts, etc.)
├── package.json        # Project dependencies and scripts
└── svelte.config.js    # SvelteKit configuration
```

## 4. Architecture

### 4.1. Frontend Architecture

The frontend is a client-side rendered SPA. SvelteKit's file-based router is used to define the application's routes. The main application logic resides within the `src/routes/(app)` directory, which is a route group for all authenticated pages.

### 4.2. State Management

Global application state is managed using Svelte's native `writable` stores, located in `src/lib/stores/index.ts`. This centralized approach provides a single source of truth for the application's data, including:

- User session information
- Application configuration
- UI state (e.g., modal visibility)
- Core application data (models, chats, etc.)

### 4.3. Backend Communication

All communication with the backend is handled by a dedicated API layer in `src/lib/apis`. This layer abstracts the details of making API calls and managing real-time connections. It utilizes a combination of:

- **REST:** For standard CRUD operations.
- **Socket.IO:** For real-time, bidirectional communication.
- **Server-Sent Events (SSE):** For unidirectional, real-time updates from the server.

## 5. Routing

The application's routing is defined by the directory structure within `src/routes`. Key routes include:

- `/`: The main application entry point.
- `/auth`: The authentication page.
- `/(app)/`: A route group for all authenticated pages, including:
    - `/c`: Chat interface
    - `/notes`: Notes management
    - `/workspace`: A general-purpose workspace

## 6. UI Components

The UI is built from a library of reusable Svelte components located in `src/lib/components`. These components are organized by feature or functionality. A detailed component reference is outside the scope of this document, but developers should familiarize themselves with the existing components before creating new ones.

## 7. Internationalization (i18n)

The application supports multiple languages using the **i18next** library. Configuration and translation files are located in `src/lib/i18n`. Components should use the `t` function provided by `i18next` to display translated strings.

## 8. Client-side Technologies

### 8.1. Pyodide

Pyodide is used to execute Python code directly in the browser. This allows for powerful client-side data processing and analysis without requiring a separate Python backend.

### 8.2. ONNX Runtime

ONNX Runtime is used to run machine learning models in the browser. This enables features like client-side image recognition, natural language processing, and other AI-powered functionality.

## 9. Getting Started

1. **Install dependencies:** `npm install`
2. **Run the development server:** `npm run dev`
3. **Build for production:** `npm run build`

## 10. Conventions

- **File Naming:** Use kebab-case for file names (e.g., `my-component.svelte`).
- **Component Naming:** Use PascalCase for Svelte component names (e.g., `<MyComponent />`).
- **Styling:** Use Tailwind CSS utility classes for styling. Avoid writing custom CSS whenever possible.
- **State Management:** Use the centralized Svelte stores for managing global state. Avoid creating local state for data that needs to be shared across components.
- **API Communication:** Use the functions provided by the `src/lib/apis` layer for all backend communication.

This document provides a high-level overview of the WebUI-Go frontend. For more detailed information, please refer to the source code and the official documentation for the libraries and frameworks used.
