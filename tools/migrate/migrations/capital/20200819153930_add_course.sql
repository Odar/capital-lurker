-- +goose Up
create table course
(
    id serial not null,
    name varchar not null,
    theme_id bigint not null,
    description text not null,
    added_at timestamp not null,
    updated_at timestamp not null,
    constraint course_fk_theme
        foreign key(theme_id)
        references theme(id)
        on delete cascade
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
