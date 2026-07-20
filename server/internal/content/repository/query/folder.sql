-- name: CreateFolder :one
insert into folders
    (workspace_id, parent_id, name, position, created_by)
values
    ($1, $2, $3, $4, $5)
returning *;

-- name: GetFolderByID :one
select * from folders where id = $1;

-- name: GetFoldersByWorkspace :many
select * from folders where workspace_id = $1
order by parent_id nulls first, position, created_at;

-- name: GetMaxPositionInParent :one
select coalesce(max(position), -1)::int as max_position
from folders
where workspace_id = $1 and parent_id is not distinct from $2;

-- name: RenameFolder :one
update folders set name = $2, updated_at = now() where id = $1 returning *;

-- name: MoveFolder :exec
update folders set
    parent_id = $2,
    position = $3,
    updated_at = now()
where id = $1;

-- name: DeleteFolder :exec
delete from folders where id = $1;

-- name: CreateDefaultFolder :one
insert into folders
    (workspace_id, parent_id, name, position, created_by, is_default)
values
    ($1, null, $2, 0, $3, true)
returning *;

-- name: LockWorkspaceStructure :exec
select pg_advisory_xact_lock(hashtext(sqlc.arg(workspace_id)::uuid::text));

-- name: ReindexFolderSiblings :exec
with ordered as (
    select f.id as folder_id,
        (row_number() over (
            order by position,
                    case when f.id = sqlc.arg(moved_id) then 0 else 1 end,
                    f.name
        ))::int - 1 as rn
    from folders f
    where f.workspace_id = sqlc.arg(workspace_id)
    and f.parent_id is not distinct from sqlc.arg(parent_id)
)
update folders t 
set position = o.rn 
from ordered o 
where t.id = o.id and t.position <> o.rn;

-- name: GetFolderByNameInParent :one
select * from folders
where workspace_id = $1
    and parent_id is not distinct from $2
    and name = $3;