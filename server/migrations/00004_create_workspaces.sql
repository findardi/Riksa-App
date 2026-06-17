-- +goose Up
-- +goose StatementBegin
create table if not exists workspaces(
    id uuid primary key default gen_random_uuid(),
    owner_id uuid not null references users (id) on delete restrict,
    name text not null,
    slug text not null,
    description text,
    status text not null default 'active',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint workspaces_owner_slug_key unique (owner_id, slug),
    constraint workspaces_status_check check (status in ('prepare', 'active', 'archive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workspaces;
-- +goose StatementEnd
