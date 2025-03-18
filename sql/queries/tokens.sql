-- name: CreateToken :one
insert into refresh_tokens (token, created_at, updated_at, user_id)
values ($1, now(), now(), $2)
returning *;

-- name: RemoveAllTokens :exec
delete from refresh_tokens;