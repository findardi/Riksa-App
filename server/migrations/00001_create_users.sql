-- +goose Up
-- +goose StatementBegin
create table if not exists users(
    id  uuid primary key default gen_random_uuid(),
    email text not null,
    username text,
    password_hash text,
    email_verified_at timestamptz,
    status text not null default 'pending',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint user_email_key unique (email),
    constraint user_status_check check (status in ('pending', 'active', 'blocked'))
);

create unique index users_username_key on users (username) where username is not null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd

