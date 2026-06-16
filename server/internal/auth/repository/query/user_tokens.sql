-- name: GetValidUserToken :one
select * from user_tokens
where user_id = $1
    and type = $2
    and used_at is null
    and expires_at > now()
order by created_at desc
limit 1;

-- name: CreateUserToken :one
insert into user_tokens 
    (user_id, type, code_hash, expires_at)
values
    ($1, $2, $3, $4)
returning *;

-- name: UpdateUserToken :one
update user_tokens set
    code_hash = $2,
    expires_at = $3,
    used_at = null
where
    user_id = $1 and type = $4
returning *;

-- name: DeleteUserToken :exec
delete from user_tokens
where code_hash = $1 and user_id = $2 and type = $3;

-- name: GetRefreshToken :one
select * from user_tokens
where code_hash = $1
 and type = 'refresh'
limit 1;

-- name: MarkRefreshTokenUsed :exec
update user_tokens set used_at = now()
where id = $1 and type = 'refresh';

-- name: DeleteTokensByType :exec
delete from user_tokens
where user_id = $1 and type = $2;

-- name: DeleteExpiredUserTokens :exec
delete from user_tokens
where user_id = $1 and expires_at < now();

-- name: GetTokenByCodeAndUser :one
select * from user_tokens
where code_hash = $1
 and type = $2
 and user_id = $3
 and used_at is null
 and expires_at > now()
limit 1;