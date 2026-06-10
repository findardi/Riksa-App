-- name: GetValidUserToken :one
select * from user_tokens
where user_id = $1
    and type = $2
    and used_at is null
    and expires_at > now()
order by created_at desc
limit 1;