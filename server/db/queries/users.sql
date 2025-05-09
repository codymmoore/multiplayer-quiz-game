-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE
    (id = $1 OR $1 IS NULL) AND
    (username = $2 OR $2 IS NULL) AND
    (email = $3 OR $3 IS NULL) AND
    is_active = true
LIMIT 1;

-- name: GetUsers :many
SELECT *
FROM users
WHERE is_active = true
LIMIT $1
OFFSET $2;

-- name: CountUsers :one
SELECT
    COUNT(*)
FROM
    users
WHERE
    is_active = true;

-- name: UpdateUser :one
UPDATE users
    SET username = CASE WHEN $2 IS NULL THEN username ELSE $2 END,
        email = CASE WHEN $3 IS NULL THEN email ELSE $3 END,
        password_hash = CASE WHEN $4 IS NULL THEN password_hash ELSE $4 END
WHERE id = $1
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users
    SET is_active = false
WHERE id = $1;