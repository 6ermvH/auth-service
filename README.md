# Auth Service

A simple and lightweight microservice for handling JWT-based authentication. This service provides endpoints for generating and refreshing tokens, using a PostgreSQL database to securely store refresh tokens.

## 🚀 Getting Started

The recommended way to run this service is with Docker Compose.

### Prerequisites

*   Docker and Docker Compose
*   A `.env` file in the project root (you can copy `.env.example` if it exists, or create one from scratch).

### 1. Configure Your Environment

Create a `.env` file in the project root with the following variables:

```env
# For JWT signing
JWT_SECRET_KEY=your_super_secret_key

# PostgreSQL connection details
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
```

You can generate a secure `JWT_SECRET_KEY` with the following command:
```bash
openssl rand -base64 64
```

### 2. Run the Service

With your `.env` file configured, start the service using Docker Compose:

```bash
docker-compose up --build -d
```

The service will be available at `http://localhost:8080`.

## ⚙️ API Endpoints

### Generate Tokens

*   **Request:**
    ```bash
    curl -X POST "http://localhost:8080/token?user_id=<your_user_uuid>"
    ```

*   **Success Response (`200 OK`):**
    ```json
    {
      "access_token": "your_jwt_access_token",
      "refresh_token": "your_base64_refresh_token"
    }
    ```

### Refresh Tokens

*   **Request:**
    ```bash
    curl -X POST http://localhost:8080/refresh \
         -H "Content-Type: application/json" \
         -d '{"access_token": "your_jwt_access_token", "refresh_token": "your_base64_refresh_token"}'
    ```

*   **Success Response (`200 OK`):**
    ```json
    {
      "access_token": "new_jwt_access_token",
      "refresh_token": "new_base64_refresh_token"
    }
    ```

## 🗄️ Database Schema

The service uses a single table to store refresh token data. The `docker-compose.yml` is configured to automatically run any `.sql` scripts placed in the `/migrations` directory upon database creation.

```sql
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    access_token_sha256 TEXT NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    client_ip TEXT NOT NULL,
    expired_at TIMESTAMP DEFAULT now() + interval \'7 days\',
    is_used BOOLEAN DEFAULT false
);
```
