INSERT INTO users (full_name, username, password, email, phone_number, address, role_status, active_role)
VALUES
  ('John Doe', 'johndoe', 'password123', 'john.doe@example.com', '1234567890', '123 Main St', 'Active', 'Customer'),
  ('Jane Smith', 'janesmith', 'password123', 'jane.smith@example.com', '1234567891', '456 Elm St', 'Active', 'Admin'),
  ('Alice Johnson', 'alicejohnson', 'password123', 'alice.johnson@example.com', '1234567892', '789 Pine St', 'Active', 'Delivery driver'),
  ('Bob Brown', 'bobbrown', 'password123', 'bob.brown@example.com', '1234567893', '101 Maple St', 'Inactive', 'Delivery driver'),
  ('Charlie Davis', 'charliedavis', 'password123', 'charlie.davis@example.com', '1234567894', '202 Oak St', 'Active', 'Customer'),
  ('David Wilson', 'davidwilson', 'password123', 'david.wilson@example.com', '1234567895', '303 Birch St', 'Active', 'Admin'),
  ('Eva Martinez', 'evamartinez', 'password123', 'eva.martinez@example.com', '1234567896', '404 Cedar St', 'Active', 'Customer'),
  ('Frank Harris', 'frankharris', 'password123', 'frank.harris@example.com', '1234567897', '505 Walnut St', 'Inactive', 'Delivery driver'),
  ('Grace Lee', 'gracelee', 'password123', 'grace.lee@example.com', '1234567898', '606 Pine St', 'Active', 'Customer'),
  ('Henry Walker', 'henrywalker', 'password123', 'henry.walker@example.com', '1234567899', '707 Oak St', 'Active', 'Admin');



INSERT INTO user_roles (user_user_id, role_role_id)
VALUES
  (1, 1), (1, 3),  
  (2, 1), (2, 3),  
  (3, 1),          
  (4, 2),       
  (5, 1),          
  (6, 1), (6, 3),  
  (7, 1),          
  (8, 2),          
  (9, 1),
  (10, 1), (10, 3);