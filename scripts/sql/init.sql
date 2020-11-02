create table if not exists mail (
    id serial primary key,
    email varchar(254) not null unique
);

create table if not exists subscription (
    email_id int references mail,
    ad_id int references advertisement

)

create table if not exists advertisement (
    id serial primary key,
    ad bigint not null unique,
    price bigint not null,
    last_check_at timestamp not null default now()
)