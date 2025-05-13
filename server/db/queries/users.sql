-- name: CreateUser :one
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE
    (id = sqlc.narg(id) OR sqlc.narg(id) IS NULL) AND
    (username = sqlc.narg(username) OR sqlc.narg(username) IS NULL) AND
    (email = sqlc.narg(email) OR sqlc.narg(email) IS NULL) AND
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
SET username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    password_hash = COALESCE(sqlc.narg(password_hash), password_hash)
WHERE id = $1
RETURNING *;

-- name: DeactivateUser :exec
UPDATE users
    SET is_active = false
WHERE id = $1;