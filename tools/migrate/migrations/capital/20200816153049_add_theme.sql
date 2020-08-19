-- +goose Up
create table theme
(
    id serial not null,
    name varchar not null,
    slug varchar not null,
    on_main_page boolean not null,
    added_at timestamp not null,
    updated_at timestamp not null,
    position int not null,
    img varchar not null
);

create unique index theme_id_uindex
    on theme (id);

create unique index theme_name_uindex
    on theme (name);

alter table theme
    add constraint theme_pk
        primary key (id);

-- +goose Down
drop table if exists theme;
