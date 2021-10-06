package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"path/routes"
	"path/config"
	"github.com/joho/godotenv" 
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	//for intialiseing thr env file
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























































/* app.Get("/", func(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "ya this [ CRUD WITHOT DATABASE ] API is working",
	})
}) */

/* 	api.Get("", func(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "y",
	})
}) */
