-- +goose Up
-- +goose StatementBegin
create table if not exists user_identities(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users (id) on delete cascade,
    provider text not null,
    provider_uid text not null,
    email text,
    created_at timestamptz not null default now(),

    constraint user_identities_provider_uid_key unique (provider, provider_uid)
);

create index user_identities_user_id_idx on user_identities (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists user_identities;
-- +goose StatementEnd
