-- +goose Up
create table video
(
    id serial not null,
    name varchar not null,
    img varchar not null,
    video varchar not null,
    youtube_video text not null,
    course_id bigint not null,
    position_in_course int not null,
    added_at timestamp not null,
    updated_at timestamp not null,
    uploaded_at timestamp not null,
    youtubed_at timestamp not null,
    constraint video_fk_course
        foreign key (course_id)
            references course (id)

);

create unique index video_id_uindex
    on video (id);

create unique index video_name_uindex
    on video (name);

alter table video
    add constraint video_pk
        primary key (id);

-- +goose Down
drop table if exists video;
