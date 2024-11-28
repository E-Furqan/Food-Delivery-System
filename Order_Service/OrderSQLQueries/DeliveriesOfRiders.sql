select 
    delivery_driver, count(*) as DeliversCompleted
from 
    orders
where
    order_status = 'Completed'
group by 
    delivery_driver
order by 
    DeliversCompleted DESC;
