    CREATE TABLE IF NOT EXISTS sources (
        sources_id SERIAL PRIMARY KEY,
        sources_name VARCHAR(100),
        source_type  VARCHAR(100)
    );