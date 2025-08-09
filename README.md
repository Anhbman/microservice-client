# Microservice Client

This is the client-side web application for the microservice project.

## Setup

1.  **Clone the repository**
    ```bash
    git clone <repository-url>
    cd microservice-client
    ```

2.  **Create an environment file**

    Copy the example environment file and update the variables as needed.
    ```bash
    cp .env.example .env
    ```

    **Environment Variables:**
    *   `SERVICE_ENDPOINT`: The URL of the backend service (e.g., `http://localhost:8081`).
    *   `PATH_TO_UPLOAD`: The path to the directory for file uploads.

3.  **Run the application**
    ```bash
    go run ./cmd/client/main.go
    ```