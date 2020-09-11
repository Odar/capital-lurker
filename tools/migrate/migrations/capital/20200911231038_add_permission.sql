-- +goose Up
create table permission
(
    id serial not null,
    name varchar(20)
);

create unique index permission_id_uindex
    on permission (id);

create unique index permission_name_uindex
    on permission (name);

-- +goose Down
drop table permission if exists;
