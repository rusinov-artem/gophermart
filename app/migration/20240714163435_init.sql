-- +goose Up
-- +goose StatementBegin
create table "user"
(
    login         text not null primary key,
    password_hash text not null
);

create table "auth_token"
(
    login text not null REFERENCES "user" (login),
    token text not null,
    PRIMARY KEY (token, login)
);

alter table auth_token
    add constraint auth_token_uniq unique (token);

create table "order"
(
    login     text           not null REFERENCES "user" (login),
    order_nr  text           not null,
    status    text           not null default 'NEW',
    upload_at timestamptz(0) not null default now(),
    accrual   bigint,
    PRIMARY KEY (order_nr, login)
);

alter table "order"
    add constraint order_uniq unique (order_nr);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "user";
drop table if exists "auth_token"
-- +goose StatementEnd
