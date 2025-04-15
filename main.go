package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/CodeMonkMI/todo/src/database"
	"github.com/CodeMonkMI/todo/src/todo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

const (
	hostname       = "mongodb://localhost:27017"
	dbURI          = "mongodb://localhost:27017"
	dbName         = "demo_todo"
	collectionName = "todo"
	port           = ":4000"
)

func main() {

	database.ConnectData()
	// disconnect from mongodb
	defer database.DisconnectData()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	rnd = renderer.New()

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
	r.Mount("/todo", todo.TodoHandlers())

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
