# 数据库设计

## Info Table

user

user_id(auto increment)|name|password|mail|account_state|join_time|role|score|gmt_create|gmt_modified|belongTo



## Submit Table

submit

submit_id(auto increment)|user_id|challenge_id|submit_time|submit_flag|correct|gmt_create|gmt_modified



## Challenge Table

challenge

challenge_id(auto increment)|description|first_blood|second_blood|third_blood|state|is_dynamic|gmt_create|gmt_modified



## Flag Table

replica_alloc

replica_alloc_id(auto increment)|replica_id|user_id|gmt_create|gmt_modified



replica

replica_id(auto increment)|gmt_create|gmt_modified|challenge_id|flag



## Bulletin Table

bulletin

bulletin_id(auto increment)|gmt_create|gmt_modified|content|topping|



