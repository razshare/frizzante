-- migration: down
drop table if exists migrations;

-- migration: up
create table if not exists migrations(
    version varchar(36) primary key,
    created_at datetime not null
);