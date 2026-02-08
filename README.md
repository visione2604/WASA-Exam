# WASA-Exam Chat Application

This is an instant messaging application developed for the Web and Software Architecture (WASA) course. It allows users to communicate in real-time through direct messages and group chats, share photos, react to messages, and manage their profiles.

## Features

*   **User Authentication**: Login and registration system.
*   **Profile Management**: Users can update their username and profile photo.
*   **Direct Messaging**: 1-on-1 private conversations.
*   **Group Chats**: Create groups, add/remove members, update group name and photo.
*   **Messaging**: Send text messages and photos.
*   **Message Interactions**:
    *   **Forward** messages to other conversations and to your own conversation.
    *   **Reply** to specific messages .
    *   **React** to messages with emojis (like, heart, laugh, sad, angry).
    *   **Delete** messages.
*   **Status Indicators**: Message delivery status (sent,read).

## Tech Stack

*   **Backend**: Go (Golang)
    *   `httprouter`
    *   `SQLite`
*   **Frontend**: Vue.js 3
    *   `Vite`
    *   `Axios`
*   **Containerization**: Docker & Docker Compose.

## Prerequisites

*   [Docker](https://www.docker.com/products/docker-desktop) installed on your machine.
*   [Git](https://git-scm.com/) (optional, to clone the repo).

## How to Run

The easiest way to run the application is using Docker Compose, which handles both the backend and frontend services.

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/visione2604/WASA-Exam.git
    cd WASA-Exam
    ```

2.  **Start the application**:
    ```bash
    docker compose up --build
    ```
    This command builds the Docker images for the backend and frontend and starts the containers.

3.  **Access the application**:
    *   **Web UI**: Open your browser and go to [http://localhost:8080](http://localhost:8080/).
    *   **API Backend**: Running at [http://localhost:3000](http://localhost:3000) (requests proxied via Nginx).

4.  **Stop the application**:

    Run:
    ```bash
    docker compose down
    ```

## API Documentation

The REST API specification is available in the `doc` directory:

*   [doc/api.yaml](doc/api.yaml): OpenAPI 3.0.3 definition describing all available endpoints, request/response bodies, and schemas.

## Project Structure

*   `cmd/webapi`: Entry point for the Go backend server.
*   `service/api`: API handlers and routing logic.
*   `service/database`: Database interactions and SQL queries.
*   `webui`: Vue.js frontend source code.
*   `doc`: Documentation and API specifications.
