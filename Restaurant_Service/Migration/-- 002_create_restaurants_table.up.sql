CREATE TABLE IF NOT EXISTS restaurants (
    RestaurantId SERIAL PRIMARY KEY,
    RestaurantName VARCHAR(255) NOT NULL,
    RestaurantAddress TEXT,
    RestaurantPhoneNumber VARCHAR(100) UNIQUE,
    RestaurantEmail VARCHAR(100) UNIQUE,
    Password VARCHAR(100),
    RestaurantStatus VARCHAR(100),
);
