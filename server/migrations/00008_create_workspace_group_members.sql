-- +goose Up
-- +goose StatementBegin
create table if not exists workspace_group_members(
    group_id uuid not null references workspace_groups (id) on delete cascade,
    member_id uuid not null references workspace_members (id) on delete cascade,
    created_at timestamptz not null default now(),

    primary key (group_id, member_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists workspace_group_members;
-- +goose StatementEnd
