CREATE TABLE IF NOT EXISTS configs (
    client_id VARCHAR(255) PRIMARY KEY,
    client_secret VARCHAR(255) NOT NULL,
    token_uri VARCHAR(255) NOT NULL,
    refresh_token TEXT,
    sources_id INT ,
    destinations_id INT,
    FOREIGN KEY (sources_id) REFERENCES sources(sources_id), 
    FOREIGN KEY (destinations_id) REFERENCES destinations(destinations_id),
    folder_Url TEXT
);