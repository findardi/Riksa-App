-- +goose Up
-- +goose StatementBegin
insert into access_levels (workspace_id, name, can_view, can_download)
values (null, 'none', false, false)
on conflict do nothing;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
delete from access_levels where workspace_id is null and name = 'none';
-- +goose StatementEnd
