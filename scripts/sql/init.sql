create table if not exists email (
    id serial primary key,
    email varchar(254) not null unique
);

create table if not exists subscription (
    email_id int not null,
    advertisement varchar(2048) not null
)