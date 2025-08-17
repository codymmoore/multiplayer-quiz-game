-- name: CreateVerificationCode :one
INSERT INTO verification_code (user_id, verification_code)
VALUES ($1, $2)
RETURNING *;

-- name: GetVerificationCode :one
SELECT *
FROM verification_code
WHERE user_id = $1
LIMIT 1;

-- name: UpdateVerificationCode :one
UPDATE verification_code
SET verification_code = $2
WHERE user_id = $1
RETURNING *;