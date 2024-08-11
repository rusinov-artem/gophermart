-- +goose Up
-- +goose StatementBegin
create table "user"
(
    login         text not null primary key,
    password_hash text not null
);

create table "auth_token"
(
    login text not null primary key,
    token text not null,
    PRIMARY KEY (login, token)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "user";
drop table if exists "auth_token"
-- +goose StatementEnd
