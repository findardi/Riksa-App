-- name: GetUserById :one
select * from users where id = $1 limit 1;

-- name: GetUserByEmail :one
select * from users where email = $1 limit 1;

-- name: GetUserByUsername :one
select * from users where username = $1 limit 1;

-- name: GetUsersByStatus :many
select * from users where status = $1;

-- name: GetUsersById :many
select * from users where id = any($1::uuid[]);

-- name: CreateUser :one
insert into users 
    (email, username, password_hash, status)
values 
    ($1, $2, $3, $4 )
returning *;

-- name: UpdateUser :one
update users set
    email = $2,
    username = $3,
    password_hash = $4,
    email_verified_at = $5,
    status = $6,
    updated_at = now()
where id = $1
returning *;