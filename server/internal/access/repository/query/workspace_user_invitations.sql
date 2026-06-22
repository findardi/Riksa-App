-- name: InsertWorkspaceInvitation :one
insert into workspace_user_invitations
    (workspace_id, email, role_id, user_id, invited_by, code_hash, status, expires_at)
values
    ($1, $2, $3, $4, $5, $6, $7, $8)
returning *;

-- name: GetWorkspaceInvitation :one
select * from workspace_user_invitations where id = $1;

-- name: GetWorkspaceInvitationByCodeHash :one
select * from workspace_user_invitations where code_hash = $1;

-- name: ListWorkspaceInvitations :many
select
    i.*,
    r.name as role_name,
    u.username as invited_by_username
from
    workspace_user_invitations i
left join
    workspace_roles r
        on r.id = i.role_id
left join
    users u
        on u.id = i.invited_by
where
    i.workspace_id = $1
    and i.status = $2
order by i.created_at desc;

-- name: AcceptWorkspaceInvitation :one
update workspace_user_invitations set
    status = 'accepted',
    user_id = $2,
    accepted_at = now(),
    updated_at = now()
where id = $1 and status = 'pending'
returning *;

-- name: RejectWorkspaceInvitation :one
update workspace_user_invitations set
    status = 'rejected',
    updated_at = now()
where id = $1 and status = 'pending'
returning *;

-- name: RevokeWorkspaceInvitation :one
update workspace_user_invitations set
    status = 'revoked',
    updated_at = now()
where id = $1 and status = 'pending'
returning *;
