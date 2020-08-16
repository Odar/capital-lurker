-- +goose Up
alter table speaker
add column university_id bigint default null;

-- +goose Down
drop sequence if exists university_id;