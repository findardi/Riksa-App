-- +goose Up
-- +goose StatementBegin
create table if not exists workspace_members(
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid not null references workspaces (id) on delete cascade,
    user_id uuid not null references users (id) on delete cascade,
    role_id uuid not null references workspace_roles (id) on delete restrict,
    status text not null default 'active',
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint workspace_members_user_key unique (workspace_id, user_id),
    constraint workspace_members_status_check check (status in ('invited', 'active', 'suspended'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workspace_members;
-- +goose StatementEnd
