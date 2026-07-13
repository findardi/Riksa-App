-- name: GetSystemAccessLevelByName :one
select * from access_levels
where workspace_id is null and name = $1;

-- name: ListAccessLevels :many
select * from access_levels
where workspace_id is null or workspace_id = $1
order by workspace_id nulls first, name;

-- name: GetAccessLevel :one
select * from access_levels
where id = $1 and (workspace_id is null or workspace_id = $2);