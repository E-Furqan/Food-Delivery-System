select 
    order_status, count(*) as statusFrequency
from 
    orders
group by 
    order_status
order by 
    statusFrequency DESC;