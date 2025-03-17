select 
    orders.order_id,
    orders.user_id,
    orders.restaurant_id,
    order_items.item_id,
    order_items.quantity,
    orders.total_bill,
    orders.delivery_driver,
    orders.order_status,
    orders.time as Order_time
from 
    orders
inner join 
    order_items on orders.order_id = order_items.order_id
where 
    order_status ='Cancelled'
Group by 
    orders.order_id,
    orders.user_id,
    orders.restaurant_id,
    orders.total_bill,
    order_items.item_id,
    order_items.quantity,
    orders.delivery_driver,
    orders.order_status,
    orders.time
limit 5 offset 0;