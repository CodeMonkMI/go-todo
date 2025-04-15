package todo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func Find() ([]todo, error) {
	cursor, err := todoCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		fmt.Println("Failed to fetch todos! from db")
		return nil, err
	}
	var todos []todoModel
	err2 := cursor.All(context.TODO(), &todos)
	if err2 != nil {
		fmt.Println("Failed to unpack todos!")
		return nil, err2

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
	return formatTodos, nil
}
func FindById(id string) (todo, error) {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid todo id")
		return todo{}, err
	}
	var todoData todoModel
	err2 := todoCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objectId}}).Decode(&todoData)
	if err2 != nil {
		fmt.Println("Invalid todo id")
		return todo{}, err2
	}

	formatTodo := todo{
		ID:        todoData.ID.Hex(),
		Title:     todoData.Title,
		Completed: todoData.Completed,
		CreatedAt: todoData.CreatedAt,
	}

	return formatTodo, nil
}
func DeleteById(id string) error {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid todo id")
		return err
	}
	_, err2 := FindById(id)
	if err2 != nil {
		fmt.Println("Invalid todo id")
		return err2
	}
	filter := bson.D{{Key: "_id", Value: objectId}}
	_, err3 := todoCollection.DeleteOne(context.TODO(), filter)
	if err3 != nil {
		fmt.Println("Delete fialed")
		return err3
	}

	return nil
}

func UpdateById(id string, data updateTodoModel) (todo, error) {

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("Invalid todo id")
		return todo{}, err
	}
	findTodo1, err2 := FindById(id)
	if err2 != nil {
		fmt.Println("Invalid todo id")
		return todo{}, err2
	}

	if data.Title == "" {
		data.Title = findTodo1.Title
	}

	filter := bson.D{{Key: "_id", Value: objectId}}
	update := bson.D{{Key: "$set", Value: updateTodo}}

	_, err3 := todoCollection.UpdateOne(context.TODO(), filter, update)
	if err3 != nil {
		fmt.Println("Failed to update todo")
		return todo{}, err3
	}
	findTodo, _ := FindById(id)

	return findTodo, nil
}
func Create(data todoModel) (todo, error) {

	_, err := todoCollection.InsertOne(context.TODO(), data)
	if err != nil {
		fmt.Println("Failed to create todo")
		return todo{}, err
	}
	todoId := data.ID
	findTodo, _ := FindById(todoId.Hex())

	return findTodo, nil
}
