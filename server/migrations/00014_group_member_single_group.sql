-- +goose Up
-- +goose StatementBegin
alter table workspace_group_members
    drop constraint workspace_group_members_pkey;

alter table workspace_group_members
    add constraint workspace_group_members_pkey primary key (member_id);

create index workspace_group_members_id_idx
    on workspace_group_members (group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists workspace_group_members_id_idx;

alter table workspace_group_members
    drop constraint workspace_group_members_pkey;

alter table workspace_group_members
    add constraint workspace_group_members_pkey primary key (group_id, member_id);
-- +goose StatementEnd
