# Auth-Service
## Build & Deploy
- Docker ```docker-compose up --build```
- Local ```cd internal && go build && ./auth_service```
## Environment 
- JWT_SECRET_KEY (generate ```openssl rand -base64 64```)
- DB_NAME
- DB_PASSWORD
- DB_USER
- DB_HOST (only local)
- DB_PORT (only local)
## PostgreSQL table
```
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    client_ip TEXT NOT NULL,
    expired_at TIMESTAMP DEFAULT now() + interval '7 days',
    is_used BOOLEAN DEFAULT false
);
```
## Experience
- Request: ```curl -X POST "http://127.0.0.1:8080/token?user_id='UUID"```
-  Response: ```{"access_token":"JWT_TOKEN","refresh_token":"BASE64_TOKEN"}```
