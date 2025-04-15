package todo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/CodeMonkMI/todo/src/utility"
	"github.com/go-chi/chi/v5"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func fetchTodos(w http.ResponseWriter, r *http.Request) {

	todos, err := Find()
	if err != nil {
		utility.ResponseError(w, err, "Failed to fetch todos")
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todos,
	})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var body todoModel
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		utility.ResponseError(w, err, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if body.Title == "" {
		utility.ResponseError(w, err, "Title is required")
		return
	}

	todoData, err := Create(todoModel{
		Title:     body.Title,
		Completed: false,
		ID:        bson.NewObjectID(),
		CreatedAt: time.Now(),
	})

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"data":    todoData,
	})

}

func fetchSingleTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	todoData, err := FindById(id)

	if err != nil {
		utility.ResponseError(w, err, "Invalid todo id")
		return
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": todoData,
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var body updateTodoModel
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utility.ResponseError(w, err, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	todoData, err2 := UpdateById(id, body)
	if err2 != nil {
		utility.ResponseError(w, err2, "Failed to update todo")
		return
	}

	rnd.JSON(w, http.StatusAccepted, renderer.M{
		"message": "Todo updated successfully",
		"data":    todoData,
	})
}
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := DeleteById(id)

	if err != nil {
		utility.ResponseError(w, err, "Todo deletion failed")
		return
	}

	rnd.JSON(w, http.StatusNoContent, renderer.M{
		"message": "Todo delete successfully",
		"data":    nil,
	})
}
