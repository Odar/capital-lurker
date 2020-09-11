-- +goose Up
create table user_role
(
    id serial not null,
    user_id serial not null,
    role_id serial not null
);

alter table user_role
    add constraint usrole_fk_role
        foreign key (role_id)
            references role (id);

alter table user_role
    add constraint usrole_fk_user
        foreign key (user_id)
            references users (id);

-- +goose Down
drop table user_role if exists;
