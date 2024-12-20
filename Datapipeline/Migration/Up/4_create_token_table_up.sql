CREATE TABLE IF NOT EXISTS tokens (
    token_id SERIAL PRIMARY KEY,
    access_token TEXT,
    token_type TEXT,
    refresh_token TEXT,
    expiry TIMESTAMP 
);

