-- +goose Up
create table user_role
(
    id bigserial not null,
    user_id serial not null,
    role_id serial not null
);

alter table user_role
    add constraint usrole_role_fk
        foreign key (role_id)
            references role (id)
            on delete cascade;

alter table user_role
    add constraint usrole_user_fk
        foreign key (user_id)
            references users (id)
            on delete cascade;

-- +goose Down
drop table user_role if exists;
