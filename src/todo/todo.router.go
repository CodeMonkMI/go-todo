package todo

import (
	"net/http"

	"github.com/CodeMonkMI/todo/src/database"
	"github.com/go-chi/chi/v5"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var rnd *renderer.Render
var todoCollection *mongo.Collection

func TodoHandlers() http.Handler {
	rnd = renderer.New()
	todoCollection = database.TodoCollection()
	r := chi.NewRouter()
	r.Get("/", fetchTodos)
	r.Post("/", createTodo)
	r.Get("/{id}", fetchSingleTodo)
	r.Patch("/{id}", updateTodo)
	r.Delete("/{id}", deleteTodo)
	return r
}
