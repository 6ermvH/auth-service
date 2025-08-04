INSERT INTO refresh_tokens (user_id, access_token_sha256, refresh_token_hash, client_ip, expired_at, is_used)
SELECT
    gen_random_uuid(),
    encode(sha256(random()::text::bytea), 'hex'),
    crypt(encode(gen_random_bytes(32), 'base64'), gen_salt('bf')),
    inet_client_addr() || (floor(random() * 255) + 1)::text,
    NOW() + (random() * 14 - 7) * '1 day'::interval,
    random() < 0.1 -- 10% of tokens will be marked as used
FROM generate_series(1, 10000);
