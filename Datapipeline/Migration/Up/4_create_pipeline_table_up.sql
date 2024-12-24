CREATE TABLE IF NOT EXISTS pipelines (
    pipeline_id SERIAL PRIMARY KEY,
    sources_id INT,
    destinations_id INT,
    FOREIGN KEY (sources_id) REFERENCES sources(sources_id),
    FOREIGN KEY (destinations_id) REFERENCES destinations(destinations_id)  
);

