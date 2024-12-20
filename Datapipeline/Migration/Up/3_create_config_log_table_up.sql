CREATE TABLE IF NOT EXISTS logs_configs (
    log_id INT PRIMARY KEY,
    config_id INT,
    FOREIGN KEY (log_id) REFERENCES logs(log_id),
    FOREIGN KEY (config_id) REFERENCES configs(config_id)  
);