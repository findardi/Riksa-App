-- name: InsertGroupMember :one
insert into workspace_group_members
    (group_id, member_id)
values
    ($1, $2)
returning *;

-- name: DeleteGroupMember :exec
delete from workspace_group_members where
    group_id = $1 and member_id = $2;

-- name: GetGroupMembers :many
select
    gm.*,
    u.username,
    u.email,
    r.name as role_name,
    g.name as group_name
from
    workspace_group_members gm
left join
    workspace_groups g on g.id = gm.group_id
left join 
    workspace_members m on m.id = gm.member_id
left join
    users u on u.id = m.user_id
left join 
    workspace_roles r on r.id = m.role_id
where
    gm.group_id = $1
order by gm.created_at;