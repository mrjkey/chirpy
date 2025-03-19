-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, user_id, expired_at)
values ($1, now(), now(), $2, $3)
returning *;

-- name: RemoveAllRefreshTokens :exec
delete from refresh_tokens;

-- name: GetAllRefreshTokens :many
select * from refresh_tokens
order by created_at desc;