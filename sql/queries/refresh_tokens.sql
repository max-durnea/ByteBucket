-- name: CreateRefreshToken :one

INSERT INTO refresh_tokens(token,created_at,updated_at,user_id,expires_at,revoked_at)
VALUES($1,$2,$3,$4,$5,$6)
RETURNING *;
-- name: GetUserFromRefreshToken :one

SELECT users.*
FROM refresh_tokens JOIN users ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

