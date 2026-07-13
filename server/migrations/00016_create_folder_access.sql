-- +goose Up
-- +goose StatementBegin
create table if not exists folder_access (
    folder_id uuid not null references folders (id) on delete cascade,
    group_id uuid not null references workspace_groups (id) on delete cascade,
    level_id uuid not null references access_levels (id) on delete restrict,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    primary key (folder_id, group_id)
);

create index folder_access_group_id_idx
    on folder_access (group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists folder_access;
-- +goose StatementEnd
