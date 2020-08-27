-- +goose Up
create table users
(
    id serial not null,
    vk_id integer not null,
    email varchar(50) not null,
    password varchar(50) not null,
    first_name varchar(30) not null,
    last_name varchar(30) not null,
    birth_date timestamp not null,
    signed_up_at timestamp not null,
    last_signed_in_at timestamp,
    updated_at timestamp not null,
    img varchar not null
);

create unique index user_id_uindex
    on users (id);

create unique index user_email_uindex
    on users (email);

create unique index user_vk_id_uindex
    on users (vk_id);

alter table users
    add constraint user_pk
        primary key (id);

-- +goose Down
drop table users if exists;
