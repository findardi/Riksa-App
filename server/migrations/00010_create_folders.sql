-- +goose Up
-- +goose StatementBegin
create table if not exists folders(
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid not null references workspaces (id) on delete cascade,
    parent_id uuid references folders (id) on delete cascade,
    name text not null,
    position integer not null default 0,
    created_by uuid not null references users (id) on delete restrict,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create unique index folders_name_per_parent_key
    on folders (workspace_id, parent_id, name) where parent_id is not null;

create unique index folders_name_root_key
    on folders (workspace_id, name) where parent_id is null;

create index folders_workspace_parent_idx on folders (workspace_id, parent_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists folders;
-- +goose StatementEnd
