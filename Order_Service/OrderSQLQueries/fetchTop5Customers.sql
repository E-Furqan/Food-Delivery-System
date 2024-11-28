select 
    user_id, COUNT(*) as order_frequency
from 
    orders
group by 
    user_id
order by 
    order_frequency DESC
limit 
    5; 