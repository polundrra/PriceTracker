create table if not exists mail (
    id serial primary key,
    email varchar(254) not null unique
);

create table if not exists subscription (
    email_id int not null,
    ad_id varchar(255) not null
)