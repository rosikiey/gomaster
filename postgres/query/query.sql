-- name: Gettodo :many
SELECT * FROM todos;

-- name: Gettodosinggle :one
SELECT * FROM todos where limits = 1;