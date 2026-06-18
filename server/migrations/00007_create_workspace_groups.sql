-- +goose Up
-- +goose StatementBegin
create table if not exists workspace_groups(
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid not null references workspaces (id) on delete cascade,
    name text not null,
    description text,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint workspace_groups_name_key unique(workspace_id, name)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workspace_groups;
-- +goose StatementEnd
