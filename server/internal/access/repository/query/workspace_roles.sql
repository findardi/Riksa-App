-- name: InsertRole :one
insert into workspace_roles
    (workspace_id, name, permissions, is_system)
values
    ($1, $2, $3, $4)
returning *;

-- name: GetRoles :many
select * from workspace_roles
where workspace_id = $1;

-- name: GetRole :one
select * from workspace_roles
where id = $1;