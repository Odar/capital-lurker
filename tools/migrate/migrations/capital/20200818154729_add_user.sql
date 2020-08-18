-- +goose Up
create table user
(
    id serial not null,
    email varchar(50) not null,
    password varchar(50) not null,
    first_name varchar(30) not null,
    last_name varchar(30) not null,
    birth_date timestamp not null,
    signed_up_at timestamp not null,
    last_signed_in_at timestamp not null,
    updated_at timestamp not null,
    img varchar not null
);

create unique index user_id_uindex
    on user (id);

create unique index user_email_uindex
    on user (email);

alter table user
    add constraint university_pk
        primary key (id);

-- +goose Down
drop table user if exists;
