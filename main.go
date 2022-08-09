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
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	//routes
	routes.TaskRoute(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app.Listen(":" + port)
}
