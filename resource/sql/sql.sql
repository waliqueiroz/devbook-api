CREATE DATABASE IF NOT EXISTS devbook;

USE devbook;

DROP TABLE IF EXISTS followers;

DROP TABLE IF EXISTS users;

CREATE TABLE users(
    id int auto_increment primary key,
    name varchar(50) not null,
    nick varchar(50) not null unique,
    email varchar(50) not null unique,
    password varchar(255) not null,
    created_at timestamp default current_timestamp()
) ENGINE = INNODB;

CREATE TABLE followers (
    user_id int not null,
    follower_id int not null,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    primary key (user_id, follower_id)
) ENGINE = INNODB;

CREATE TABLE publications(
    id int auto_increment primary key,
    title varchar(255) not null,
    content text not null,
    author_id int not null,
    likes int not null default 0,
    created_at timestamp default current_timestamp()
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
) ENGINE = INNODB;