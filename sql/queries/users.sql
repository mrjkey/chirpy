-- name: CreateUser :one
insert into users (id, created_at, updated_at, email, hashed_password)
values (gen_random_uuid(), now(), now(), $1, $2)
returning *;


-- name: RemoveAllUsers :exec
delete from users;

-- name: GetUserByEmail :one
select * from users
where email = $1;

-- name: UpdateUser :one
update users
set hashed_password = $2, email = $3, updated_at = now()
where id = $1
returning *;

-- name: UpgradeUserToRed :one
update users 
set is_chirpy_red = true
where id = $1
returning *;

-- name: DowngradeUser :one
update users 
set is_chirpy_red = false
where id = $1
returning *;