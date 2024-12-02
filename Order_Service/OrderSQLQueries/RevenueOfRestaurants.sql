select 
    restaurant_id , sum(total_bill) as revenue
from 
    orders
GROUP by 
    restaurant_id
Order by 
    revenue DESC;


