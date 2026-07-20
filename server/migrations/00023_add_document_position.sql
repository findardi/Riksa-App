-- +goose Up
-- +goose StatementBegin
alter table documents
    add column position integer not null default 0;

update documents d set position = s.rn
from (
    select id,
        (row_number() over (partition by folder_id order by name, created_at))::int - 1 as rn
    from documents
) s 
where d.id = s.id;

create index documents_folder_position_idx on
    documents(folder_id, position);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop index if exists documents_folder_position_idx;
alter table documents
    drop column if exists position;
-- +goose StatementEnd
