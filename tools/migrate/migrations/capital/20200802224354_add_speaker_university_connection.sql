-- +goose Up
alter table speaker
add column university_id bigint default null;