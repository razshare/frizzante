-- migration: up
create table if not exists sessions(
    id varchar(36) primary key,
    created_at datetime not null,
    updated_at datetime not null,
    error text not null default ''
);
create table if not exists todos(
    id varchar(36) primary key,
    session_id varchar(36) not null,
    description varchar(256) not null default '',
    checked int not null default 0,
    foreign key (session_id) references sessions(id)
);
-- migration: down
drop table if exists sessions;
drop table if exists todos;