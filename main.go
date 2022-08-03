package main

import (
	"fiber-mongo-api/configs"
	"fiber-mongo-api/routes" //add this
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	//routes
	routes.UserRoute(app) //add this
	routes.TaskRoute(app)
	app.Use(cors.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app.Listen(":" + port)
}
