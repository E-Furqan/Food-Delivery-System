select 
    delivery_driver, count(*) as DeliversCompleted
from 
    orders
where
    order_status = 'Completed' AND delivery_driver != 0
group by 
    delivery_driver
order by 
    DeliversCompleted DESC;
