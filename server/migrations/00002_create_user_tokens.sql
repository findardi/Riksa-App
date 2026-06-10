-- +goose Up
-- +goose StatementBegin
create table if not exists user_tokens(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users (id) on delete cascade,
    type text not null,
    code_hash text not null,
    expires_at timestamptz not null,
    used_at timestamptz,
    created_at timestamptz not null default now(),

    constraint users_token_type_check check (type in ('email_verification', 'password_reset', 'refresh'))
);

create index user_tokens_user_type_idx on user_tokens (user_id, type);
create index user_tokens_expires_idx on user_tokens (expires_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_tokens;
-- +goose StatementEnd
