select COUNT(*) as active_user_count from users where active_role='Customer' and role_status='Active';

select COUNT(*) as active_user_count from users where active_role='Delivery driver' and role_status='Active';

select COUNT(*) as active_user_count from users where active_role='Admin' and role_status='Active';  