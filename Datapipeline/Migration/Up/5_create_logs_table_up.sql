CREATE TABLE IF NOT EXISTS logs (
    log_id SERIAL PRIMARY KEY,
    log_message TEXT,
    pipelines_id INT,
    FOREIGN KEY (pipelines_id) REFERENCES pipelines(pipeline_id), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
            
);

