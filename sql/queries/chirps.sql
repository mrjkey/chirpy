-- name: AddChirp :one
insert into chirps (id, created_at, updated_at, body, user_id)
values (gen_random_uuid(), now(), now(), $1, $2)
returning *;

-- name: RemoveChirps :exec
delete from chirps;

-- name: DeleteChirp :exec
delete from chirps
where id = $1;

-- name: GetAllChirps :many
select * from chirps
order by created_at asc;

-- name: GetChirpById :one
select * from chirps
where id = $1;