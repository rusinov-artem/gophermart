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
    accrual   double precision,
    PRIMARY KEY (login, order_nr)
);

alter table "order"
    add constraint order_order_nr_uniq unique (order_nr);

create table "withdraw"
(
    login      text           not null REFERENCES "user" (login),
    order_nr   text           not null,
    created_at timestamptz(0) not null default now(),
    sum        double precision,
    PRIMARY KEY (login, order_nr)
);

alter table "withdraw"
    add constraint withdraw_order_nr_uniq unique (order_nr);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists "user";
drop table if exists "auth_token";
drop table if exists "order";
drop table if exists "withdraw";
-- +goose StatementEnd
