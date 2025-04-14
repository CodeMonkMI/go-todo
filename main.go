package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var rnd *renderer.Render
var db *mongo.Database
var todoCollection *mongo.Collection

const (
	hostname       = "mongodb://localhost:27017"
	dbURI          = "mongodb://localhost:27017"
	dbName         = "demo_todo"
	collectionName = "todo"
	port           = ":4000"
)

type (
	todoModel struct {
		ID        bson.ObjectID `bson:"_id"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"createdAt"`
	}
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

func fetchTodos(w http.ResponseWriter, r *http.Request) {

	cursor, err := todoCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println("Failed to fetch todos! from db")
		responseError(w, err, "failed to fetch todos")
		return
	}
	var todos []todoModel
	err2 := cursor.All(context.TODO(), &todos)
	if err2 != nil {
		fmt.Println("Failed to unpack todos!")
		log.Println(err2)
		responseError(w, err2, "failed to fetch todos")
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
		responseError(w, err, "Invalid request payload")
		return
	}

	defer r.Body.Close()

	if newTodo.Title == "" {
		responseError(w, err, "Title is required")
		return
	}

	newTodo.ID = bson.NewObjectID()
	newTodo.CreatedAt = time.Now()
	result, err := todoCollection.InsertOne(context.TODO(), newTodo)
	if err != nil {
		fmt.Println(err)
		responseError(w, err, "Failed to create todo")
		return
	}
	todoId := result.InsertedID.(bson.ObjectID)

	var todoData todoModel
	err3 := todoCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: todoId}}).Decode(&todoData)
	if err3 != nil {

		responseError(w, err3, "Failed to create todo")
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

func updateTodo(w http.ResponseWriter, r *http.Request) {

}

func main() {

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	fmt.Println("Starting application...")
	rnd = renderer.New()
	// connect to mongodb
	client, error := mongo.Connect(options.Client().ApplyURI(dbURI))
	checkErr(error, "failed connect to mongodb")

	// get collections
	todoCollection = client.Database(dbName).Collection(collectionName)

	// define chi routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// home handler
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		rnd.JSON(w, http.StatusOK, renderer.M{
			"message": "Server is running!",
		})
		return
	})
	r.Mount("/todo", todoHandlers())

	srv := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Println("Starting server on port", port)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan

	log.Println("Shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()

	// disconnect from mongodb
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Println("Database client disconnected")
			panic(err)
		}
	}()

}

func todoHandlers() http.Handler {
	r := chi.NewRouter()
	r.Get("/", fetchTodos)
	r.Post("/", createTodo)
	r.Put("/{id}", updateTodo)
	return r
}

func responseError(w http.ResponseWriter, err error, msg string) {
	rnd.JSON(w, http.StatusBadRequest, renderer.M{
		"message": msg,
		"err":     err,
	})
	return
}

func checkErr(e error, customMsg string) {
	if e != nil {
		fmt.Println(customMsg)
		log.Fatal(e)
	}
}
