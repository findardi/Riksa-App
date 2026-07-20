-- name: CreateDocument :one
insert into documents 
    (workspace_id, folder_id, name, position, uploaded_by)
values
    ($1, $2, $3, $4, $5)
returning *;

-- name: CreateDocumentVersion :one
insert into document_versions
    (document_id, version_no, mime, size, storage_key, uploaded_by)
values 
    ($1, $2, $3, $4, $5, $6)
returning *;

-- name: SetCurrentVersion :exec
update documents set 
    current_version_id = $2,
    updated_at = now()
where id = $1;

-- name: GetNextVersionNo :one
select coalesce(max(version_no),0)::int + 1 as next_no
from document_versions where document_id = $1;

-- name: GetDocumentByID :one
select * from documents where id = $1;

-- name: ListDocumentsByFolder :many
select
    d.id,
    d.name,
    d.folder_id,
    d.current_version_id,
    d.uploaded_by,
    d.created_at,
    d.updated_at,
    v.version_no,
    v.mime,
    v.size
from documents d
join document_versions v on v.id = d.current_version_id
where d.folder_id = $1
order by d.position, d.name, d.created_at;

-- name: ListVersionByDocument :many
select * from document_versions where document_id = $1 order by version_no desc;

-- name: GetVersionByID :one
select * from document_versions where id = $1;

-- name: GetCurrentVersion :one
select v.* from document_versions v 
join documents d on d.current_version_id = v.id
where d.id = $1;

-- name: DeleteDocument :exec
delete from documents where id = $1;

-- name: MoveDocument :exec
update documents set folder_id = $2, position = $3, updated_at = now() where id = $1;

-- name: SetVersionRendition :exec
update document_versions
set rendition_key = sqlc.arg(rendition_key),
    page_count = sqlc.arg(page_count)
where id = sqlc.arg(id);

-- name: GetMaxPosition :one
select coalesce(max(position), -1)::int as max_position
from documents
where folder_id = $1;

-- name: ReindexDocumentSiblings :exec
with ordered as (
    select d.id as document_id,
        (row_number() over (
            order by position,
                    case when d.id = sqlc.arg(moved_id) then 0 else 1 end,
                    d.name 
        ))::int - 1 as rn 
    from documents d
    where d.folder_id = sqlc.arg(folder_id)
)
update documents t
set position = o.rn 
from ordered o 
where t.id = o.document_id and t.position <> o.rn;