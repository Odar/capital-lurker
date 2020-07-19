-- +goose Up
create table receiver_cache
(
	id serial not null,
	word varchar not null,
	answer varchar not null,
	uses int not null
);

create unique index receiver_cache_id_uindex
	on receiver_cache (id);

create unique index receiver_cache_word_uindex
	on receiver_cache (word);

alter table receiver_cache
	add constraint receiver_cache_pk
		primary key (id);



-- +goose Down
DROP TABLE receiver_cache IF EXISTS;
