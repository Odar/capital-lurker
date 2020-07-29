-- +goose Up
create table speaker
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

create unique index speaker_id_uindex
	on speaker (id);

create unique index speaker_name_uindex
    on speaker (name);

alter table speaker
	add constraint speaker_pk
		primary key (id);

-- +goose Down
DROP TABLE speaker IF EXISTS;
