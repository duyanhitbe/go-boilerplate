-- name: GetAllTodo :many
SELECT * FROM "todos";

-- name: CreateTodo :one
INSERT INTO "todos" ("title")
VALUES ($1)
RETURNING *;

-- name: DeleteAllTodo :exec
DELETE FROM "todos";