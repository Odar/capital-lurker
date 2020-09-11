-- +goose Up
create table role
(
    id serial not null,
    name varchar(20)
);

create unique index role_id_uindex
    on role (id);

create unique index role_name_uindex
    on role (name);

-- +goose Down
drop table role if exists;
