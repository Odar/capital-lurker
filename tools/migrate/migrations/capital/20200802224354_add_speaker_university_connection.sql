-- +goose Up
alter table speaker
add column university_id bigint default null;

alter table speaker
    add constraint speaker_university_fk
        foreign key (university_id)
            references university (id)
            on delete set null;