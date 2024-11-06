CREATE TABLE IF NOT EXISTS items (
    ItemId SERIAL PRIMARY KEY,
    ItemName VARCHAR(255) NOT NULL,
    ItemDescription TEXT,
    ItemPrice NUMERIC(10, 2) NOT NULL,
    RestaurantId INTEGER REFERENCES restaurants(RestaurantId)
);
