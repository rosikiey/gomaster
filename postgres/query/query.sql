-- name: Gettodo :many
SELECT * FROM todos;

-- name: Gettodosinggle :one
SELECT * FROM todos where id = $1 limit 1;

-- name: CreateTodo :one
insert INTO TODOS (NAME) values ($1) returning *;

-- name: DeleteTodobyId :exec
DELETE FROM TODOS WHERE ID = $1;

-- name: UpdateTodo :one
UPDATE TODOS SET NAME = $2, COMPLETED = $3 WHERE ID = $1 returning *;