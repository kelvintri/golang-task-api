package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes" //add this
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(app) //add this
	routes.TaskRoute(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app.Listen(":" + port)
}
