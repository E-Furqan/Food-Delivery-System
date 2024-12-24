CREATE TABLE IF NOT EXISTS configs (
    config_id SERIAL PRIMARY KEY,
    client_id VARCHAR(255),
    client_secret VARCHAR(255) NOT NULL,
    token_uri VARCHAR(255) NOT NULL,
    refresh_token TEXT,
    sources_id INT ,
    destinations_id INT,
    folder_Url TEXT
);