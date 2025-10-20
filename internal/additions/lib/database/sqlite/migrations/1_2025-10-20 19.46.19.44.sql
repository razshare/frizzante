-- migration: down
drop table user;

-- migration: up
create table user(
    id varchar(36) primary key
);
