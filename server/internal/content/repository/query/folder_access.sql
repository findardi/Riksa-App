-- name: SetFolderAccess :one
insert into folder_access (folder_id, group_id, level_id)
select f.id, g.id, l.id
from folders f
join workspace_groups g
    on g.id = sqlc.arg(group_id) and g.workspace_id = f.workspace_id
join access_levels l
    on l.id = sqlc.arg(level_id) and (l.workspace_id is null or l.workspace_id = f.workspace_id)
where f.id = sqlc.arg(folder_id) and f.workspace_id = sqlc.arg(workspace_id)
on conflict (folder_id, group_id) do update
    set
        level_id = excluded.level_id,
        updated_at = now()
returning *;

-- name: RemoveFolderAccess :exec
delete from folder_access fa 
using folders f 
where fa.folder_id = f.id
and fa.folder_id = $1
and fa.group_id = $2
and f.workspace_id = $3;

-- name: ListFolderAccess :many
select
    fa.folder_id,
    fa.group_id,
    g.name as group_name,
    fa.level_id,
    l.name as level_name,
    l.can_view,
    l.can_download,
    l.can_watermark
from folder_access fa
join folders f on f.id = fa.folder_id
join workspace_groups g on g.id = fa.group_id
join access_levels l on l.id = fa.level_id
where fa.folder_id = $1 and f.workspace_id = $2
order by g.name;