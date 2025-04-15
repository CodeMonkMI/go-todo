package todo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type (
	todoModel struct {
		ID        bson.ObjectID `bson:"_id"`
		Title     string        `bson:"title"`
		Completed bool          `bson:"completed"`
		CreatedAt time.Time     `bson:"createdAt"`
	}
	updateTodoModel struct {
		Title     string `bson:"title"`
		Completed bool   `bson:"completed"`
	}
	todo struct {
		ID        string    `json:"id"`
		Title     string    `json:"title"`
		Completed bool      `json:"completed"`
		CreatedAt time.Time `json:"createdAt"`
	}
)
