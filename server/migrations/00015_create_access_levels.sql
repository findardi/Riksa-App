-- +goose Up
-- +goose StatementBegin
create table if not exists access_levels (
    id uuid primary key default gen_random_uuid(),
    workspace_id uuid references workspaces (id) on delete cascade,
    name text not null,

    can_view boolean not null default false,
    can_download boolean not null default false,
    can_watermark boolean not null default false,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    constraint access_levels_download_needs_view check (not can_download or can_view),
    constraint access_levels_watermark_needs_view check (not can_watermark or can_view)
);

create unique index access_levels_system_name_key
    on access_levels (name) where workspace_id is null;

create unique index access_levels_workspace_name_key
    on access_levels (workspace_id, name) where workspace_id is not null;

insert into access_levels (workspace_id, name, can_view, can_download)
values
    (null, 'view',  true, false),
    (null, 'download', true, true);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists access_levels;
-- +goose StatementEnd
