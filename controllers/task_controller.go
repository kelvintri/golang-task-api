package controllers

import (
	"context"
	"fiber-mongo-api/configs"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var taskCollection *mongo.Collection = configs.GetCollection(configs.DB, "tasks")

func CreateTask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var task models.Task
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TaskResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&task); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TaskResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	newTask := models.Task{
		Id:          primitive.NewObjectID(),
		TaskName:    task.TaskName,
		Assignee:    task.Assignee,
		TaskDone:    task.TaskDone,
		CreatedDate: time.Now(),
		DeadLine:    task.DeadLine,
	}

	result, err := taskCollection.InsertOne(ctx, newTask)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.TaskResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetATask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskId := c.Params("taskId")
	var task models.Task
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)

	err := taskCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&task)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.TaskResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": task}})
}

func EditATask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskId := c.Params("taskId")
	var task models.Task
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)

	//validate the request body
	if err := c.BodyParser(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TaskResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&task); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.TaskResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	update := bson.M{"taskname": task.TaskName, "assignee": task.Assignee, "taskdone": task.TaskDone, "deadline": task.DeadLine}

	result, err := taskCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated task details
	var updatedTask models.Task
	if result.MatchedCount == 1 {
		err := taskCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedTask)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	return c.Status(http.StatusOK).JSON(responses.TaskResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedTask}})
}

func DeleteATask(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskId := c.Params("taskId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)

	result, err := taskCollection.DeleteOne(ctx, bson.M{"id": objId})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.TaskResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Task with specified ID not found!"}},
		)
	}

	return c.Status(http.StatusOK).JSON(
		responses.TaskResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Task successfully deleted!"}},
	)
}

func GetAllTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var tasks []models.Task
	defer cancel()

	results, err := taskCollection.Find(ctx, bson.M{})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleTask models.Task
		if err = results.Decode(&singleTask); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		tasks = append(tasks, singleTask)
	}

	return c.Status(http.StatusOK).JSON(
		responses.TaskResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": tasks}},
	)
}

func MarkATaskAsDone(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	taskId := c.Params("taskId")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(taskId)

	result, err := taskCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": bson.M{"taskdone": true}})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if result.MatchedCount == 1 {
		return c.Status(http.StatusOK).JSON(
			responses.TaskResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "Task successfully marked as done!"}},
		)
	}

	return c.Status(http.StatusNotFound).JSON(
		responses.TaskResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "Task with specified ID not found!"}},
	)
}

//list of task that done
func GetAllDoneTasks(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var tasks []models.Task
	defer cancel()

	results, err := taskCollection.Find(ctx, bson.M{"taskdone": true})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//reading from the db in an optimal way
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleTask models.Task
		if err = results.Decode(&singleTask); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.TaskResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		tasks = append(tasks, singleTask)
	}

	return c.Status(http.StatusOK).JSON(
		responses.TaskResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": tasks}},
	)
}
