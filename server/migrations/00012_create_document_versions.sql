-- +goose Up
-- +goose StatementBegin
create table if not exists document_versions(
    id uuid primary key default gen_random_uuid(),
    document_id uuid not null references documents (id) on delete cascade,
    version_no integer not null,
    mime text not null,
    size bigint not null,
    storage_key text not null,
    uploaded_by uuid not null references users (id) on delete restrict,
    created_at timestamptz not null default now(),

    constraint document_versions_no_key unique (document_id, version_no)
);

alter table documents
    add constraint documents_current_version_fk
    foreign key (current_version_id) references document_versions (id) on delete set null;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
alter table documents drop constraint if exists documents_current_version_fk;
drop table if exists document_versions;
-- +goose StatementEnd
