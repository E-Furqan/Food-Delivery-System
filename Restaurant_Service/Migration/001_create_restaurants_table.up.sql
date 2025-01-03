CREATE TABLE IF NOT EXISTS restaurants (
    restaurant_id SERIAL PRIMARY KEY,
    restaurant_name VARCHAR(255) NOT NULL,
    restaurant_address TEXT,
    restaurant_phone_number VARCHAR(15) UNIQUE,
    restaurant_status  VARCHAR(100) NOT NULL,
    restaurant_email VARCHAR(100) UNIQUE,
    password VARCHAR(100)
);
