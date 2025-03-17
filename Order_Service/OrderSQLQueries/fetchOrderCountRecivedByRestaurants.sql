SELECT 
    restaurant_id, COUNT(*) AS completed_orders
FROM 
    orders
WHERE 
    order_status = 'Completed' 
    AND time BETWEEN '2024-11-28 10:46:38.27784+00' and '2024-12-30 10:54:28.715338+00' 
GROUP BY 
    restaurant_id;