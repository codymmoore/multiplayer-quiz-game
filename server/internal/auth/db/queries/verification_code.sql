-- name: UpsertVerificationCode :one
INSERT INTO verification_code (user_id, verification_code_hash)
VALUES ($1, $2)
ON CONFLICT (user_id)
DO UPDATE SET
    verification_code_hash = EXCLUDED.verification_code_hash,
    expires_at = CURRENT_TIMESTAMP + INTERVAL '15 minutes'
RETURNING *;

-- name: GetVerificationCode :one
SELECT *
FROM verification_code
WHERE user_id = $1 AND
      expires_at > CURRENT_TIMESTAMP
LIMIT 1;