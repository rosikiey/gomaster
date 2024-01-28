package util

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/rosikiey/gomaster.git/main"
	"github.com/rosikiey/gomaster.git/postgres"
)

func mapTodo(todo postgres.Todo) interface{} {
	return struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Completed bool   `json:"completed"`
	}{
		ID:        todo.ID,
		Name:      todo.Name,
		Completed: todo.Completed.Bool,
	}
}

func SetupTodosRoutes(grp fiber.Router, handlers *main.Handlers) {
	todosRoutes := grp.Group("/todos")
	todosRoutes.Get("/", handlers.GetTodos)
	todosRoutes.Post("/", handlers.CreateTodo)
	todosRoutes.Get("/:id", handlers.GetTodo)
	todosRoutes.Delete("/:id", handlers.DeleteTodo)
	todosRoutes.Patch("/:id", handlers.UpdateTodo)
}

func (h *main.Handlers) UpdateTodo(ctx *fiber.Ctx) {
	type request struct {
		Name      *string `json:"name"`
		Completed *bool   `json:"completed"`
	}

	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	var body request
	err = ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse body",
		})
		return
	}

	todo, err := h.Repo.Gettodosinggle(ctx.Context(), int64(id))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	if body.Name != nil {
		todo.Name = *body.Name
	}

	if body.Completed != nil {
		todo.Completed = sql.NullBool{
			Bool:  *body.Completed,
			Valid: true,
		}
	}

	todo, err = h.Repo.UpdateTodo(ctx.Context(), postgres.UpdateTodoParams{
		ID:        int64(id),
		Name:      todo.Name,
		Completed: todo.Completed,
	})
	if err != nil {
		ctx.SendStatus(fiber.StatusNotFound)
		return
	}

	if err := ctx.Status(fiber.StatusOK).JSON(mapTodo(todo)); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}
}

func (h *main.Handlers) DeleteTodo(ctx *fiber.Ctx) {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	_, err = h.Repo.Gettodosinggle(ctx.Context(), int64(id))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	err = h.Repo.DeleteTodobyId(ctx.Context(), int64(id))
	if err != nil {
		ctx.SendStatus(fiber.StatusNotFound)
		return
	}

	ctx.SendStatus(fiber.StatusNoContent)
}

func (h *main.Handlers) GetTodo(ctx *fiber.Ctx) {
	paramsId := ctx.Params("id")
	id, err := strconv.Atoi(paramsId)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	todo, err := h.Repo.Gettodosinggle(ctx.Context(), int64(id))
	if err != nil {
		ctx.Status(fiber.StatusNotFound)
		return
	}

	if err := ctx.Status(fiber.StatusOK).JSON(mapTodo(todo)); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}
}

func (h *main.Handlers) CreateTodo(ctx *fiber.Ctx) {
	type request struct {
		Name string `json:"name"`
	}

	var body request

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse json",
		})
		return
	}

	if len(body.Name) <= 2 {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name not long enough",
		})
		return
	}

	todo, err := h.Repo.CreateTodo(ctx.Context(), body.Name)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}

	if err := ctx.Status(fiber.StatusCreated).JSON(mapTodo(todo)); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}
}

func (h *main.Handlers) GetTodos(ctx *fiber.Ctx) {
	todos, err := h.Repo.Gettodo(ctx.Context())
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}

	result := make([]interface{}, len(todos))
	for i, todo := range todos {
		result[i] = mapTodo(todo)
	}

	if err := ctx.Status(fiber.StatusOK).JSON(result); err != nil {
		ctx.Status(fiber.StatusInternalServerError).Send(err.Error())
		return
	}
}
