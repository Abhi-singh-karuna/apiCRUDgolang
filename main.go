package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"log"
	"path/config"
	"path/routes"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	//for initialize the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config.ConnectDB()

	setRoutes(app)

	err = app.Listen(":9000")
	if err != nil {
		log.Fatal(err)
	}

}

func setRoutes(app *fiber.App) {

	api := app.Group("/")

	routes.BookRutes(api.Group("/book"))

}
