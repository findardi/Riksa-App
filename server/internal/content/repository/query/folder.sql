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