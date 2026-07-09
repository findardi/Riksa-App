-- +goose Up
-- +goose StatementBegin
alter table folders add column is_default boolean not null default false;

create unique index folders_one_default_per_workspace
    on folders (workspace_id) where is_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists folders_one_default_per_workspace;
alter table folders drop column if exists is_default;
-- +goose StatementEnd
