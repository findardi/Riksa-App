-- +goose Up
-- +goose StatementBegin
alter table document_versions
    add column if not exists rendition_key text,
    add column if not exists page_count integer;

alter table document_versions
    add constraint document_versions_rendition_pair_check
    check ((rendition_key is null) = (page_count is null));

alter table document_versions
    add constraint document_versions_page_count_check
    check (page_count is null or page_count > 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table document_versions drop constraint if exists document_versions_page_count_check;
alter table document_versions drop constraint if exists document_versions_rendition_pair_check;
alter table document_versions drop column if exists rendition_key;
alter table document_versions drop column if exists page_count;
-- +goose StatementEnd
