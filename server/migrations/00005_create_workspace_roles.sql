-- +goose Up
-- +goose StatementBegin
create table if not exists workspace_roles(
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid not null references workspaces (id) on delete cascade,
    name text not null,
    permissions text[] not null default '{}',
    is_system boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint workspace_roles_name_key unique (workspace_id, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workspace_roles;
-- +goose StatementEnd
