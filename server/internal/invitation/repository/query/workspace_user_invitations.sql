-- name: GetMyInvitations :many
select 
    wi.*,
    u.username as invited_name,
    r.name as role_name,
    w.name as workspace_name
from 
    workspace_user_invitations wi
left join
    workspaces w on w.id = wi.workspace_id
left join
    users u on u.id = wi.invited_by
left join
    workspace_roles r on r.id = wi.role_id
where wi.user_id = $1 and wi.status = 'pending' and wi.expires_at > now()
order by wi.created_at desc;

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

-- name: GetWorkspaceInvitation :one
select * from workspace_user_invitations where id = $1;

-- name: InsertMember :exec
insert into workspace_members
    (workspace_id, user_id, role_id, status)
values
    ($1, $2, $3, 'active')
on conflict (workspace_id, user_id) do nothing;