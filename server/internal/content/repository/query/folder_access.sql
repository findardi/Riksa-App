-- name: SetFolderAccess :one
insert into folder_access (folder_id, group_id, can_view, can_download, can_watermark)
select f.id, g.id, sqlc.arg(can_view), sqlc.arg(can_download), sqlc.arg(can_watermark)
from folders f
join workspace_groups g
    on g.id = sqlc.arg(group_id) and g.workspace_id = f.workspace_id
where f.id = sqlc.arg(folder_id) and f.workspace_id = sqlc.arg(workspace_id)
on conflict (folder_id, group_id) do update
    set
        can_view = excluded.can_view,
        can_download = excluded.can_download,
        can_watermark = excluded.can_watermark,
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
    fa.can_view,
    fa.can_download,
    fa.can_watermark
from folder_access fa
join folders f on f.id = fa.folder_id
join workspace_groups g on g.id = fa.group_id
where fa.folder_id = $1 and f.workspace_id = $2
order by g.name;

-- name: ResolveFolderAccess :one
with recursive chain as (
    select f.id, f.parent_id, 0 as depth
    from folders f
    where f.id = sqlc.arg(folder_id) and f.workspace_id = sqlc.arg(workspace_id)

    union all

    select f.id, f.parent_id, c.depth + 1
    from folders f
    join chain c on f.id = c.parent_id
)
select
    fa.can_view,
    fa.can_download,
    fa.can_watermark
from chain c
join folder_access fa on fa.folder_id = c.id
join workspace_group_members gm on gm.group_id = fa.group_id
join workspace_members m on m.id = gm.member_id
where m.workspace_id = sqlc.arg(workspace_id) and m.user_id = sqlc.arg(user_id)
order by c.depth
limit 1;

-- name: ListVisibleFolders :many
with recursive granted as (
    select
        f.id,
        f.parent_id,
        f.name,
        f.position,
        f.is_default,
        fa.can_view,
        fa.can_download
    from folders f
    join folder_access fa on fa.folder_id = f.id
    join workspace_group_members gm on gm.group_id = fa.group_id
    join workspace_members m on m.id = gm.member_id
    where f.workspace_id = sqlc.arg(workspace_id)
      and m.workspace_id = sqlc.arg(workspace_id)
      and m.user_id = sqlc.arg(user_id)

    union all

    select
        c.id,
        c.parent_id,
        c.name,
        c.position,
        c.is_default,
        g.can_view,
        g.can_download
    from folders c
    join granted g on c.parent_id = g.id
    where not exists (
        select 1
        from folder_access fa2
        join workspace_group_members gm2 on gm2.group_id = fa2.group_id
        join workspace_members m2 on m2.id = gm2.member_id
        where fa2.folder_id = c.id
          and m2.workspace_id = sqlc.arg(workspace_id)
          and m2.user_id = sqlc.arg(user_id)
    )
)
select id, parent_id, name, position, is_default, can_view, can_download
from granted
where can_view
order by position;
