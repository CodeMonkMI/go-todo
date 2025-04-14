package main

// import (
// 	"context"
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"string"
// 	"time"

// 	"github.com/go-chi/chi"
// 	"github.com/go-chi/chi/middleware"
// 	"github.com/thedevsaddam/renderer"
// 	mgo "gopkg.in/mgo.v2"
// 	"gopkg.in/mgo.v2/bson"
// )
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var rnd *renderer.Render
var db *mgo.Database

const (
	hostname       = "mongodb://localhost:27017"
	dbName         = "demo_todo"
	collectionName = "todo"
	port           = "4000"
)

type (
	todoModel struct {
		ID        string    `bson:"id,omitempty"`
		Title     string    `bson:"title"`
		Completed bool      `bson:"completed"`
		CreatedAt time.Time `bson:"createdAt"`
	}
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"createdAt"`
	}
)

func init() {
	fmt.Println("Starting application...")
	rnd = renderer.New()
	session, mongoErr := mgo.Dial(hostname)
	checkErr(mongoErr, "Failed to connect to mongo database ")
	session.SetMode(mgo.Monotonic, true)

}

func fetchTodos(w http.ResponseWriter, r *http.Request) {
	todos := []todoModel{}
	if err := db.C(collectionName).Find(bson.M{}).All((&todos)); err != nil {
		rnd.JSON(w, http.StatusProcessing, renderer.M{
			"message": "Failed to fetch todos!",
			"err":     err,
		})
		return
	}
	todoList := []todo{}
	for _, t := range todos {
		todoList = append(todoList, todo{
			ID:        t.ID,
			Title:     t.Title,
			Completed: t.Completed,
		})
	}
	rnd.JSON(w, http.StatusOK, renderer.M{
		"message": "successfully fetched todos!",
		"data":    todoList,
	})
}
func createTodo(w http.ResponseWriter, r *http.Request) {

}
func updateTodo(w http.ResponseWriter, r *http.Request) {

}

func main() {
	fmt.Println("Application is in main function")
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

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

}

func todoHandlers() http.Handler {
	r := chi.NewRouter()
	r.Get("/", fetchTodos)
	r.Post("/", createTodo)
	r.Put("/{id}", updateTodo)
	return r
}

func checkErr(e error, customMsg string) {
	if e != nil {
		fmt.Println(customMsg)
		log.Fatal(e)
	}
}
