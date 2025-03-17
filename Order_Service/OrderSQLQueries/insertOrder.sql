-- Inserting Orders
INSERT INTO orders (user_id, restaurant_id, delivery_driver, order_status, total_bill, time)
VALUES
  -- Order by John Doe (UserID = 1) at The Gourmet House (RestaurantID = 1)
  (1, 1, 3, 'Completed', 68.95, CURRENT_TIMESTAMP),
  
  -- Order by Charlie Davis (UserID = 5) at Pizza Paradise (RestaurantID = 2)
  (5, 2, 4, 'Pending', 52.95, CURRENT_TIMESTAMP),
  
  -- Order by Eva Martinez (UserID = 7) at Sushi World (RestaurantID = 3)
  (7, 3, 3, 'Pending', 59.95, CURRENT_TIMESTAMP),
  
  -- Order by Grace Lee (UserID = 10) at Taco Fiesta (RestaurantID = 4)
  (10, 4, 4, 'Completed', 29.95, CURRENT_TIMESTAMP);





INSERT INTO items (item_id)
VALUES
    (1),
    (2),    
    (3),
    (4),
    (5),
    (6),
    (7),
    (9),
    (10),
    (11),
    (13),
    (14),
    (15);


-- Inserting Items for the Orders
INSERT INTO order_items (order_id, item_id, quantity)
VALUES
  -- Items for John Doe's Order (OrderID = 1)
  (1, 1, 1),  -- Gourmet Burger
  (1, 2, 2),  -- Truffle Fries
  (1, 3, 1),  -- Lobster Roll
  
  -- Items for Charlie Davis' Order (OrderID = 2)
  (2, 5, 2),  -- Margherita Pizza
  (2, 6, 1),  -- Pepperoni Pizza
  (2, 7, 1),  -- BBQ Chicken Pizza
  
  -- Items for Eva Martinez' Order (OrderID = 3)
  (3, 9, 2),  -- California Roll
  (3, 10, 1), -- Tuna Sashimi
  (3, 11, 1), -- Salmon Nigiri
  
  -- Items for Grace Lee's Order (OrderID = 4)
  (4, 13, 3), -- Beef Taco
  (4, 14, 1), -- Chicken Quesadilla
  (4, 15, 2); -- Fish Taco




INSERT INTO orders (user_id, restaurant_id, delivery_driver, order_status, total_bill, time)
VALUES
  (3, 1, 2, 'Completed', 72.45, CURRENT_TIMESTAMP),
  (6, 2, 3, 'Pending', 60.80, CURRENT_TIMESTAMP),
  (8, 3, 4, 'Completed', 84.30, CURRENT_TIMESTAMP),
  (9, 4, 1, 'Completed', 56.25, CURRENT_TIMESTAMP),
  (10, 1, 4, 'Pending', 47.99, CURRENT_TIMESTAMP),
  (2, 1, 3, 'Completed', 82.90, CURRENT_TIMESTAMP),
  (4, 2, 4, 'Pending', 75.50, CURRENT_TIMESTAMP),
  (6, 3, 2, 'Completed', 95.20, CURRENT_TIMESTAMP),
  (9, 4, 3, 'Completed', 68.75, CURRENT_TIMESTAMP),
  (7, 1, 1, 'Pending', 50.60, CURRENT_TIMESTAMP),
  (8, 2, 2, 'Completed', 83.00, CURRENT_TIMESTAMP),
  (10, 3, 4, 'Pending', 91.45, CURRENT_TIMESTAMP),
  (3, 4, 1, 'Completed', 52.40, CURRENT_TIMESTAMP),
  (5, 1, 2, 'Pending', 66.35, CURRENT_TIMESTAMP),
  (11, 2, 3, 'Completed', 55.30, CURRENT_TIMESTAMP),
  (1, 4, 2, 'Pending', 56.30, CURRENT_TIMESTAMP),
  (2, 3, 3, 'Completed', 64.90, CURRENT_TIMESTAMP),
  (3, 2, 1, 'Completed', 72.50, CURRENT_TIMESTAMP),
  (4, 1, 4, 'Pending', 88.75, CURRENT_TIMESTAMP),
  (5, 4, 2, 'Completed', 47.60, CURRENT_TIMESTAMP),
  (6, 3, 3, 'Completed', 98.20, CURRENT_TIMESTAMP),
  (7, 2, 1, 'Pending', 55.45, CURRENT_TIMESTAMP),
  (8, 1, 4, 'Completed', 81.90, CURRENT_TIMESTAMP),
  (9, 3, 2, 'Pending', 90.60, CURRENT_TIMESTAMP),
  (10, 4, 1, 'Completed', 63.75, CURRENT_TIMESTAMP);


INSERT INTO order_items (order_id, item_id, quantity)
VALUES
  (5, 1, 1),
  (5, 2, 2),
  (5, 3, 1),
  (6, 5, 1),
  (6, 6, 2),
  (6, 7, 1),
  (7, 9, 2),
  (7, 10, 1),
  (7, 11, 1),
  (8, 13, 2),
  (8, 14, 1),
  (8, 15, 1),
  (9, 1, 2),
  (9, 2, 1),
  (9, 3, 2),
  (10, 5, 1),
  (10, 6, 1),
  (10, 7, 2),
  (11, 9, 1),
  (11, 10, 2),
  (11, 11, 1),
  (12, 13, 2),
  (12, 14, 1),
  (12, 15, 1),
  (13, 9, 1),
  (13, 10, 2),
  (13, 11, 1),
  (14, 5, 1),
  (14, 6, 2),
  (14, 7, 1),
  (15, 1, 3),
  (15, 2, 2),
  (15, 3, 1),
  (16, 9, 2),
  (16, 10, 1),
  (16, 11, 2),
  (17, 5, 1),
  (17, 6, 1),
  (17, 7, 1),
  (18, 1, 2),
  (18, 2, 2),
  (18, 3, 1),
  (19, 9, 2),
  (19, 10, 1),
  (19, 11, 1),
  (20, 13, 1),
  (20, 14, 2),
  (20, 15, 1),
  (23, 14, 2),
  (23, 15, 1),
  (22, 13, 1),
  (22, 14, 2),
  (22, 15, 1);





INSERT INTO orders (user_id, restaurant_id, delivery_driver, order_status, total_bill, time)
VALUES
  (1, 1, 2, 'Cancelled', 45.50, CURRENT_TIMESTAMP),
  (2, 3, 4, 'Cancelled', 67.30, CURRENT_TIMESTAMP),
  (3, 2, 1, 'Cancelled', 72.00, CURRENT_TIMESTAMP),
  (4, 4, 3, 'Cancelled', 58.90, CURRENT_TIMESTAMP),
  (1, 1, 2, 'Cancelled', 39.75, CURRENT_TIMESTAMP),
  (2, 3, 4, 'Cancelled', 61.40, CURRENT_TIMESTAMP),
  (3, 2, 3, 'Cancelled', 70.15, CURRENT_TIMESTAMP),
  (4, 4, 1, 'Cancelled', 49.30, CURRENT_TIMESTAMP),
  (1, 1, 2, 'Cancelled', 52.80, CURRENT_TIMESTAMP),
  (2, 3, 4, 'Cancelled', 78.25, CURRENT_TIMESTAMP);



INSERT INTO order_items (order_id, item_id, quantity)
VALUES
  (30, 1, 1),
  (30, 2, 2),
  (31, 3, 1),
  (31, 5, 1),
  (32, 6, 2),
  (32, 7, 1),
  (33, 9, 2),
  (33, 1, 2),
  (35, 10, 1),
  (35, 11, 1),
  (34, 10, 2),
  (34, 14, 1),
  (36, 15, 1),
  (36, 11, 1),
  (37, 10, 2),
  (37, 14, 1),
  (38, 15, 1),
  (38, 1, 1),
  (39, 11, 1),
  (39, 10, 2);