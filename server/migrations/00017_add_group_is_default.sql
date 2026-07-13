-- +goose Up
-- +goose StatementBegin
alter table workspace_groups add column is_default boolean not null default false;

create unique index workspace_groups_one_default_per_workspace
    on workspace_groups (workspace_id) where is_default;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists workspace_groups_one_default_per_workspace;
alter table workspace_groups drop column if exists is_default;
-- +goose StatementEnd
