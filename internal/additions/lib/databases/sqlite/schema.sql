-- migration: down
drop table if exists user;

-- migration: up
create table user(
    id varchar(36) primary key
);