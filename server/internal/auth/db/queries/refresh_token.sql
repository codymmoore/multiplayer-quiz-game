-- name: CreateRefreshToken :one
INSERT INTO refresh_token (token_hash, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: GetRefreshToken :one
SELECT *
FROM refresh_token
WHERE token_hash = $1 AND
      active = true AND
      expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: DeactivateRefreshToken :exec
UPDATE refresh_token
SET active = false
WHERE id = $1;

-- name: DeactivateUserRefreshTokens :exec
UPDATE refresh_token
SET active = false
WHERE user_id = $1;