package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func TaskRoute(app *fiber.App) {
	app.Post("/task", controllers.CreateTask)
	app.Get("/task/:taskId", controllers.GetATask)
	app.Put("/task/:taskId", controllers.EditATask)
	app.Delete("/task/:taskId", controllers.DeleteATask)
	app.Put("/task/:taskId/done", controllers.MarkATaskAsDone)
	app.Get("/tasks", controllers.GetAllTasks)
}
