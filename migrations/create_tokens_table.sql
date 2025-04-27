CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    refresh_token_hash TEXT NOT NULL,
    client_ip TEXT NOT NULL,
    expired_at TIMESTAMP DEFAULT now() + interval '7 days',
    is_used BOOLEAN DEFAULT false
);

INSERT INTO refresh_tokens (user_id, refresh_token_hash, client_ip)
VALUES
('11111111-1111-1111-1111-111111111111', crypt('test_refresh_hash_1', gen_salt('bf')), '127.0.0.1'),
('22222222-2222-2222-2222-222222222222', crypt('test_refresh_hash_2', gen_salt('bf')), '127.0.0.1');

