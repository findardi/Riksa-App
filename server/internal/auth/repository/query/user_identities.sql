-- name: GetUserIdentity :one
select * from user_identities
where provider = $1
    and provider_uid = $2
limit 1;