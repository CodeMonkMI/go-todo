package todo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CodeMonkMI/todo/src/utility"
	"github.com/go-chi/chi/v5"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func fetchTodos(w http.ResponseWriter, r *http.Request) {

	cursor, err := todoCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println("Failed to fetch todos! from db")

		utility.ResponseError(w, err, "failed to fetch todos")
		return
	}
	var todos []todoModel
	err2 := cursor.All(context.TODO(), &todos)
	if err2 != nil {
		fmt.Println("Failed to unpack todos!")
		log.Println(err2)
		utility.ResponseError(w, err2, "failed to fetch todos")
		return

	}

	var formatTodos []todo
	for _, t := range todos {

		formatTodos = append(formatTodos, todo{
			ID:        t.ID.Hex(),
			Title:     t.Title,
			Completed: t.Completed,
			CreatedAt: t.CreatedAt,
		})
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": formatTodos,
	})
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var newTodo todoModel
	err := json.NewDecoder(r.Body).Decode(&newTodo)

	if err != nil {
		utility.ResponseError(w, err, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if newTodo.Title == "" {
		utility.ResponseError(w, err, "Title is required")
		return
	}

	newTodo.ID = bson.NewObjectID()
	newTodo.CreatedAt = time.Now()
	result, err := todoCollection.InsertOne(context.TODO(), newTodo)
	if err != nil {
		fmt.Println(err)
		utility.ResponseError(w, err, "Failed to create todo")
		return
	}
	todoId := result.InsertedID.(bson.ObjectID)

	var todoData todoModel
	err3 := todoCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: todoId}}).Decode(&todoData)
	if err3 != nil {

		utility.ResponseError(w, err3, "Failed to create todo")
		return
	}

	formatTodo := todo{
		ID:        todoData.ID.Hex(),
		Title:     todoData.Title,
		Completed: todoData.Completed,
		CreatedAt: todoData.CreatedAt,
	}

	rnd.JSON(w, http.StatusCreated, renderer.M{
		"message": "Todo created successfully",
		"data":    formatTodo,
	})

}

func fetchSingleTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utility.ResponseError(w, err, "Invalid todo id")
		return
	}
	var todoData todoModel
	err2 := todoCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}).Decode(&todoData)
	if err2 != nil {
		utility.ResponseError(w, err2, "Invalid todo id")
		return
	}
	formatTodo := todo{
		ID:        todoData.ID.Hex(),
		Title:     todoData.Title,
		Completed: todoData.Completed,
		CreatedAt: todoData.CreatedAt,
	}

	rnd.JSON(w, http.StatusOK, renderer.M{
		"data": formatTodo,
	})
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utility.ResponseError(w, err, "Invalid todo id")
		return
	}
	var findTodo todoModel
	err2 := todoCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}).Decode(&findTodo)
	if err2 != nil {
		utility.ResponseError(w, err2, "Invalid todo id")
		return
	}

	var updateTodo updateTodoModel
	err3 := json.NewDecoder(r.Body).Decode(&updateTodo)
	if err3 != nil {
		utility.ResponseError(w, err3, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if updateTodo.Title == "" {
		updateTodo.Title = findTodo.Title
	}
	if !updateTodo.Completed {
		updateTodo.Completed = findTodo.Completed
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: updateTodo}}

	_, err4 := todoCollection.UpdateOne(context.TODO(), filter, update)
	if err4 != nil {
		utility.ResponseError(w, err4, "Failed to update todo")
		return
	}

	rnd.JSON(w, http.StatusAccepted, renderer.M{
		"message": "Todo updated successfully",
		"data": todo{
			ID:        findTodo.ID.Hex(),
			Title:     updateTodo.Title,
			Completed: updateTodo.Completed,
			CreatedAt: findTodo.CreatedAt,
		},
	})
}
func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		utility.ResponseError(w, err, "Invalid todo id")
		return
	}
	var findTodo todoModel
	err2 := todoCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}).Decode(&findTodo)
	if err2 != nil {
		utility.ResponseError(w, err2, "Invalid todo id")
		return
	}

	filter := bson.D{{Key: "_id", Value: objectId}}

	_, err3 := todoCollection.DeleteOne(context.TODO(), filter)

	if err3 != nil {
		utility.ResponseError(w, err3, "Todo deletion failed")
		return
	}

	rnd.JSON(w, http.StatusNoContent, renderer.M{
		"message": "Todo delete successfully",
		"data":    nil,
	})
}
