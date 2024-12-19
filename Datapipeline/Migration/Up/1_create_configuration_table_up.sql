CREATE TABLE IF NOT EXISTS configs (
    config_id SERIAL PRIMARY KEY,
    client_id VARCHAR(255) UNIQUE NOT NULL,
    client_secret VARCHAR(255) NOT NULL,
    token_uri VARCHAR(255) NOT NULL,
    auth_uri VARCHAR(255) NOT NULL,
    redirect_uris TEXT NOT NULL,
    auth_provider_cert_url VARCHAR(255) NOT NULL
);