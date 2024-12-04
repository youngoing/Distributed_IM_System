create DATABASE backend;

use backend;

CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE Table user_details (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_id int NOT NULL,
    nickname VARCHAR(255) NOT NULL UNIQUE,
    avatar_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table chat_groups (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_detail_id int NOT NULL,
    name VARCHAR(255) NOT NULL,
    avatar_url VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

create table group_members (
    id INT PRIMARY KEY AUTO_INCREMENT,
    group_id int NOT NULL,
    user_detail_id int NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_friends (
    id INT PRIMARY KEY AUTO_INCREMENT,
    user_detail_id int NOT NULL,
    friend_id int NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 避免重复记录
CREATE UNIQUE INDEX idx_user_friend_unique ON user_friends (user_detail_id, friend_id);