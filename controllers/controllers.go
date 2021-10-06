package controllers

import (
	"fmt"
	"os"
	"path/config"
	"path/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/gofiber/fiber/v2"
)

// controller to get all the book from data slice

func GetBooks(c *fiber.Ctx) error {

	// Collection gets a handle for a collection with the given name configured with the given CollectionOptions
	BookCollection := config.MI.DB.Collection(os.Getenv("DATABASE_COLLECTION"))

	// Query to filter or iterate over the all the element
	query := bson.D{{}}

	// Context returns *fasthttp.RequestCtx that carries a deadline
	cursor, err := BookCollection.Find(c.Context(), query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	var books []models.Book = make([]models.Book, 0)

	// iterate the cursor and decode each item into a Books
	//All iterates the cursor and decodes each document into results.
	//The results parameter must be a pointer to a slice.

	err = cursor.All(c.Context(), &books)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"book": books,
		},
	})

}

// add new book into databse

func CreareBook(c *fiber.Ctx) error {
	bookCollection := config.MI.DB.Collection(os.Getenv("DATABASE_COLLECTION"))

	// The new built-in function allocates memory. The first argument is a type,
	// not a value, and the value returned is a pointer to a newly
	// allocated zero value of that type.
	data := new(models.Book)

	// BodyParser binds the request body to a struct.
	err := c.BodyParser(&data)

	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})

	}

	//shows the data from struct and this is also be an empty
	data.ID = nil
	f := false
	data.Completed = &f
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	// InsertOne executes an insert command to insert a single document into the collection.
	result, err := bookCollection.InsertOne(c.Context(), data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot insert todo",
			"error":   err,
		})
	}

	//fetch and print the inserted data
	book := &models.Book{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	//printed tha above inserted data
	bookCollection.FindOne(c.Context(), query).Decode(book)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"book": book,
		},
	})
}

// fetch a single book by its ID
func GetBook(c *fiber.Ctx) error {

	bookCollection := config.MI.DB.Collection(os.Getenv("DATABASE_COLLECTION"))

	//get parameter value
	bookid := c.Params("id")

	id, _ := primitive.ObjectIDFromHex(bookid)

	book := &models.Book{}

	query := bson.D{{Key: "_id", Value: id}}

	err := bookCollection.FindOne(c.Context(), query).Decode(book)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Todo Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"book": book,
		},
	})

}

//update in your books

func UpdateBook(c *fiber.Ctx) error {

	todoCollection := config.MI.DB.Collection(os.Getenv("DATABASE_COLLECTION"))

	// find parameter
	paramID := c.Params("id")

	// convert parameterID to objectId
	id, err := primitive.ObjectIDFromHex(paramID)

	// if parameter cannot parse
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	// var data Request
	data := new(models.Book)
	err = c.BodyParser(&data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
			"error":   err,
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	// updateData
	var dataToUpdate bson.D

	if data.Title != nil {
		// todo.Title = *data.Title
		dataToUpdate = append(dataToUpdate, bson.E{Key: "title", Value: data.Title})
	}

	if data.Completed != nil {
		// todo.Completed = *data.Completed
		dataToUpdate = append(dataToUpdate, bson.E{Key: "completed", Value: data.Completed})
	}

	dataToUpdate = append(dataToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	update := bson.D{
		{Key: "$set", Value: dataToUpdate},
	}

	// update
	err = todoCollection.FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Todo Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot update todo",
			"error":   err,
		})
	}

	// get updated data
	todo := &models.Book{}

	todoCollection.FindOne(c.Context(), query).Decode(todo)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"todo": todo,
		},
	})

}

// delete a book
func DeleteBook(c *fiber.Ctx) error {

	bookCollection := config.MI.DB.Collection(os.Getenv("DATABASE_COLLECTION"))

	//get param
	bookid := c.Params("id")

	//convert parameter to object id
	id, err := primitive.ObjectIDFromHex(bookid)

	//if parameter is not exist
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "cannot parse id",
			"error":   err,
		})
	}

	// find and delete the books from the list
	query := bson.D{{Key: "_id", Value: id}}

	err = bookCollection.FindOneAndDelete(c.Context(), query).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Todo Not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot delete todo",
			"error":   err,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)

}
