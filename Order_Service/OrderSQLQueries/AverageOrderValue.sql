select
    user_id, 
    Avg(total_bill) As Average_Order_Value
from
    orders
where 
    user_id = 2
group by 
    user_id;



select
    restaurant_id, 
    Avg(total_bill) As Average_Order_Value
from
    orders
where 
    restaurant_id = 4
group by 
    restaurant_id;


select
    time, 
    Avg(total_bill) As Average_Order_Value
from
    orders
where 
    time between '2024-11-28 10:53:38.27784+00' and '2024-11-28 11:58:38.27784+00'
group by 
    time;