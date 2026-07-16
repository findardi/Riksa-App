-- +goose Up
-- +goose StatementBegin
alter table folder_access
    add column can_download_original boolean not null default false,
    add column can_share boolean not null default false;

alter table folder_access
    add constraint folder_access_download_orignal_needs_download
        check (not can_download_original or can_download),
    add constraint folder_access_share_need_view
        check (not can_share or can_view);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table folder_access
    drop constraint folder_access_download_orignal_needs_download,
    drop constraint folder_access_share_need_view,
    drop column can_download_original,
    drop column can_share;
-- +goose StatementEnd
