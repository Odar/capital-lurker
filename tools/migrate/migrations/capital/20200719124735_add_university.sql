-- +goose Up
create table university
(
    id serial not null,
    name varchar not null,
    on_main_page boolean not null,
    in_filter boolean not null,
    added_at timestamp not null,
    updated_at timestamp not null,
    position int not null,
    img varchar not null
);

create unique index univesity_id_uindex
    on university (id);

create unique index university_name_uindex
    on university (name);

alter table university
    add constraint university_pk
        primary key (id);

-- +goose Down
drop table if exists university;
