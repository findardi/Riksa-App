-- +goose Up
-- +goose StatementBegin
alter table folder_access
    add column can_view boolean not null default false,
    add column can_download boolean not null default false,
    add column can_watermark boolean not null default false;

update folder_access fa
set can_view = l.can_view,
    can_download = l.can_download,
    can_watermark = l.can_watermark
from access_levels l
where l.id = fa.level_id;

alter table folder_access
    drop column level_id;

alter table folder_access
    add constraint folder_access_download_needs_view
        check (not can_download or can_view),
    add constraint folder_access_watermark_needs_view
        check (not can_watermark or can_view);

drop table access_levels;
-- +goose StatementEnd

-- +goose Down
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

insert into access_levels (workspace_id, name, can_view, can_download, can_watermark)
values
    (null, 'none', false, false, false),
    (null, 'view', true, false, true),
    (null, 'download', true, true, false)
on conflict do nothing;

alter table folder_access
    add column level_id uuid references access_levels (id) on delete restrict;

update folder_access fa
set level_id = l.id
from access_levels l
where l.workspace_id is null
  and l.name = case
        when fa.can_download then 'download'
        when fa.can_view then 'view'
        else 'none'
      end;

alter table folder_access
    alter column level_id set not null,
    drop constraint folder_access_watermark_needs_view,
    drop constraint folder_access_download_needs_view,
    drop column can_watermark,
    drop column can_download,
    drop column can_view;
-- +goose StatementEnd
