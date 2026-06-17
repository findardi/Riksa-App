-- name: CreateWorkspace :one
insert into workspaces 
    (owner_id, name, slug, description, status)
values 
    ($1, $2, $3, $4, $5)
returning *;

-- name: GetWorkspacesByOwner :many
select * from workspaces where owner_id = $1;

-- name: GetWorkspaceBySlugAndOwner :one
select * from workspaces 
where owner_id = $1 and slug = $2;

-- name: GetWorkspaceByID :one
select * from workspaces where id = $1;

-- name: UpdateWorkspaceStatus :exec
update workspaces set 
    status = $2,
    updated_at = now()
where id = $1;

-- name: UpdateWorkspace :one
update workspaces set
    name = $2,
    slug = $3,
    description = $4,
    updated_at = now()
where id = $1
returning *;

-- name: DeleteWorkspace :exec
delete from workspaces where id = $1;