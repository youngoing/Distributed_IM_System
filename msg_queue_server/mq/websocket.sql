CREATE DATABASE websocket;
use websocket;
create Table offline_private_message(
    id int primary key auto_increment,
    receiver_id VARCHAR(100) not null,
    msg JSON not null,
    create_time timestamp not null default current_timestamp
); 

create Table offline_group_message(
    id int primary key auto_increment,
    receiver_id VARCHAR(100) not null,
    group_id VARCHAR(100) not null,
    msg JSON not null,
    create_time timestamp not null default current_timestamp
); 