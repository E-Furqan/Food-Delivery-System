INSERT INTO restaurants (restaurant_name, restaurant_address, restaurant_phone_number, restaurant_status, restaurant_email, password)
VALUES
  ('The Gourmet House', '123 Culinary St', '1234567890', 'Open', 'gourmet@house.com', 'password1'),
  ('Pizza Paradise', '456 Cheese Ave', '9876543210', 'Open', 'pizza@paradise.com', 'password2'),
  ('Sushi World', '789 Fish Ln', '5551234567', 'Closed', 'sushi@world.com', 'password3'),
  ('Taco Fiesta', '101 Taco Blvd', '5559876543', 'Open', 'taco@fiesta.com', 'password4');



INSERT INTO items (item_name, item_description, item_price, restaurant_id)
VALUES
  -- Dishes for The Gourmet House (restaurant_id = 1)
  ('Gourmet Burger', 'A delicious gourmet beef burger with all the fixings.', 15.99, 1),
  ('Truffle Fries', 'Crispy fries with truffle oil and parmesan.', 8.99, 1),
  ('Lobster Roll', 'Fresh lobster in a buttery roll.', 22.99, 1),
  ('Grilled Salmon', 'Freshly grilled salmon with a side of vegetables.', 18.99, 1),
  
  -- Dishes for Pizza Paradise (restaurant_id = 2)
  ('Margherita Pizza', 'Classic pizza with mozzarella, tomatoes, and basil.', 12.99, 2),
  ('Pepperoni Pizza', 'Pepperoni slices with mozzarella cheese.', 14.99, 2),
  ('BBQ Chicken Pizza', 'BBQ chicken, red onions, and mozzarella cheese.', 16.99, 2),
  ('Veggie Supreme', 'A pizza loaded with fresh veggies and mozzarella.', 13.99, 2),
  ('Hawaiian Pizza', 'Ham, pineapple, and mozzarella on a pizza base.', 15.99, 2),
  
  -- Dishes for Sushi World (restaurant_id = 3)
  ('California Roll', 'Sushi rolls with crab, avocado, and cucumber.', 9.99, 3),
  ('Tuna Sashimi', 'Fresh tuna slices served with soy sauce.', 18.99, 3),
  ('Salmon Nigiri', 'Fresh salmon on top of vinegared rice.', 10.99, 3),
  ('Spicy Tuna Roll', 'Tuna with spicy mayo and avocado in a roll.', 12.99, 3),
  
  -- Dishes for Taco Fiesta (restaurant_id = 4)
  ('Beef Taco', 'Taco with seasoned beef, lettuce, and cheese.', 3.99, 4),
  ('Chicken Quesadilla', 'Grilled tortilla with chicken and cheese filling.', 8.99, 4),
  ('Fish Taco', 'Fresh fish fillet with cabbage and crema.', 4.99, 4),
  ('Churros', 'Fried dough with cinnamon sugar.', 5.99, 4);
