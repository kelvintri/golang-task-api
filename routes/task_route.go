package routes

import (
	"fiber-mongo-api/controllers"

	"github.com/gofiber/fiber/v2"
)

func TaskRoute(app *fiber.App) {
	app.Post("/task", controllers.CreateTask)
	app.Get("/task/:userId", controllers.GetATask)
	app.Put("/task/:userId", controllers.EditATask)
	app.Delete("/task/:userId", controllers.DeleteATask)
	app.Get("/tasks", controllers.GetAllTasks)
}
