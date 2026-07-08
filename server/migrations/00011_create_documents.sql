-- +goose Up
-- +goose StatementBegin
create table if not exists documents(
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid not null references workspaces (id) on delete cascade,
    folder_id uuid not null references folders (id) on delete cascade,
    name text not null,
    current_version_id uuid,
    uploaded_by uuid not null references users (id) on delete restrict,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);

create index documents_folders_idx on documents (folder_id, created_at);
create index documents_workspace_idx on documents (workspace_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists documents;
-- +goose StatementEnd
