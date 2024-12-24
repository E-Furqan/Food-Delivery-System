    CREATE TABLE IF NOT EXISTS destinations (
        destinations_id SERIAL PRIMARY KEY,
        destinations_name VARCHAR(100),
        storage_type  VARCHAR(100)
    );