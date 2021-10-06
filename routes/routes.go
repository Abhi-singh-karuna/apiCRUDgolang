package routes

import (
	"path/controllers"

	"github.com/gofiber/fiber/v2"
)

func BookRutes(route fiber.Router) {
	route.Get("", controllers.GetBooks)
	route.Post("", controllers.CreareBook)
	route.Get("/:id", controllers.GetBook)
	route.Put("/:id", controllers.UpdateBook)
	route.Delete("/:id", controllers.DeleteBook)
}
