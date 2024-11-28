
--day
select
    TO_CHAR(time, 'HH12:MI') AS hours_of_date,
    count(total_bill) As Total_Orders
from
    orders
where 
    time between '2024-11-28' and '2024-11-29'
group by 
    hours_of_date
order by 
    Total_Orders DESC;




--week
select
    TO_CHAR(time, 'DD') AS date_of_week,
    count(total_bill) As Total_Orders
from
    orders
where 
    time between '2024-11-28' and '2024-12-5'
group by 
    date_of_week
order by 
    Total_Orders DESC;



--month
SELECT
    EXTRACT(WEEK FROM time) - EXTRACT(WEEK FROM DATE_TRUNC('month', time)) + 1 AS week_of_month,
    COUNT(total_bill) AS Total_Orders
FROM
    orders
WHERE 
    time BETWEEN '2024-11-28' AND '2024-12-29'
GROUP BY 
    week_of_month
ORDER BY 
    Total_Orders DESC;




--year
select
    TO_CHAR(time, 'MM') AS month_of_year,
    count(total_bill) As Total_Orders
from
    orders
where 
    time between '2024-11-28' and '2025-11-29'
group by 
    month_of_year
order by 
    Total_Orders DESC;

