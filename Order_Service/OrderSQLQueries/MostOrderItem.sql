SELECT 
    order_items.item_id, 
    COUNT(*) AS purchaseCount, 
    orders.restaurant_id
FROM 
    order_items
INNER JOIN 
    orders ON orders.order_id = order_items.order_id
GROUP BY 
    order_items.item_id, orders.restaurant_id
ORDER BY 
    purchaseCount DESC
LIMIT 5;
