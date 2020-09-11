-- +goose Up
create table permission_role
(
    id serial not null,
    role_id serial not null,
    permission_id serial not null
);

alter table permission_role
    add constraint permrole_fk_role
        foreign key (role_id)
        references role (id);

alter table permission_role
    add constraint permrole_fk_permission
        foreign key (permission_id)
            references permission (id);

-- +goose Down
drop table permission_role if exists;
