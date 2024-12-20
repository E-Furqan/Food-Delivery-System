    CREATE TABLE IF NOT EXISTS tokens_configs (
        token_id INT PRIMARY KEY,
        config_id INT,
        FOREIGN KEY (token_id) REFERENCES tokens(token_id),
        FOREIGN KEY (config_id) REFERENCES configs(config_id)  
    );