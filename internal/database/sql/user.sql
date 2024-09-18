-- name: FetchUser :many
SELECT *
FROM users
LIMIT $1
OFFSET $2;

-- name: FindOneUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING *;