-- +goose Up
create table course
(
    id serial not null,
    name varchar not null,
    theme_id int not null,
    description text not null,
    position int not null,
    added_at timestamp not null,
    updated_at timestamp not null,
    constraint course_theme_fk
        foreign key (theme_id)
        references theme (id)
);

create unique index course_id_uindex
    on course (id);

create unique index course_name_uindex
    on course (name);

alter table course
    add constraint course_pk
        primary key (id);

-- +goose Down
drop table if exists course;
