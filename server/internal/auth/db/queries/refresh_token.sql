-- name: CreateRefreshToken :one
INSERT INTO refresh_token (token_hash, issued_at, expires_at, active)
VALUES ($1, $2, $3, true)
RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_token
WHERE token_hash = $1
LIMIT 1;

-- name: DeactivateRefreshToken :one
UPDATE refresh_token
SET active = false
WHERE id = $1
RETURNING *;