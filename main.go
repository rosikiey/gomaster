package main

import (
	"database/sql"

	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"

	"github.com/rosikiey/gomaster.git/postgres"
	"github.com/rosikiey/gomaster.git/util"

	_ "github.com/lib/pq"
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

type Handlers struct {
	Repo *postgres.Repo
}

func NewHandlers(repo *postgres.Repo) *Handlers {
	return &Handlers{Repo: repo}
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:Veg@zr01@103.127.99.34:5432/gomaster?sslmode=disable")

	if err != nil {
		panic(err)
	}

	repo := postgres.NewRepo(db)

	app := fiber.New()

	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("hello world")
	})

	handlers := NewHandlers(repo)

	SetupApiV1(app, handlers)

	err = app.Listen(3000)
	if err != nil {
		panic(err)
	}
}

func SetupApiV1(app *fiber.App, handlers *Handlers) {
	v1 := app.Group("/v1")

	util.SetupTodosRoutes(v1, handlers)
}
