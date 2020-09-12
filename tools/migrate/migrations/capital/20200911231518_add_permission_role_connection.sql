-- +goose Up
create table permission_role
(
    id bigserial not null,
    role_id serial not null,
    permission_id serial not null
);

alter table permission_role
    add constraint permrole_role_fk
        foreign key (role_id)
            references role (id)
            on delete cascade;

alter table permission_role
    add constraint permrole_permission_fk
        foreign key (permission_id)
            references permission (id)
            on delete cascade;

-- +goose Down
drop table permission_role if exists;
