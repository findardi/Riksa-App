-- name: AddMember :one
insert into workspace_members
    (workspace_id, user_id, role_id, status)
values
    ($1, $2, $3, $4)
returning *;