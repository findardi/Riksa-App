-- name: GetUserIdentity :one
select * from user_identities
where provider = $1
    and provider_uid = $2
limit 1;

-- name: CreateUserIdentity :one
insert into user_identities (user_id, provider, provider_uid, email)
values ($1, $2, $3, $4)
returning *;