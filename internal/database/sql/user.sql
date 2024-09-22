-- name: FindOneUserByUsername :one
SELECT *
FROM users
WHERE username = $1;

-- name: FindOneUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (username, password)
VALUES ($1, $2)
RETURNING *;