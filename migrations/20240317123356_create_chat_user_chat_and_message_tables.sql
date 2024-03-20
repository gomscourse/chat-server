-- +goose Up
create table chat
(
    id         bigserial primary key,
    created_at timestamp not null default now(),
    updated_at timestamp
);

create table user_chat
(
    id         bigserial primary key,
    chat_id    bigint references chat (id) not null,
    username   varchar(50)                 not null,
    created_at timestamp                   not null default now(),
    updated_at timestamp
);

create table message
(
    id         bigserial primary key,
    chat_id    bigint references chat (id) not null,
    author       varchar(50)                 not null,
    created_at timestamp                   not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table chat;
drop table user_chat;
drop table message;

