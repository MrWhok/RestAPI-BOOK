package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MrWhok/go-fiber-postgres/models"
	"github.com/MrWhok/go-fiber-postgres/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"` //what is json:.. means? it will convert the struct json author:"data" to struct book Author:"data"
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api") //anyone who start with api. have /api
	api.Post("/create_books", r.CreateBook)
	api.Delete("/delete_book/:id", r.DeleteBook)
	api.Get("/get_books", r.GetBooks)
	api.Get("/get_book/:id", r.GetBookByID)
}

func (r *Repository) CreateBook(context *fiber.Ctx) error { //why using error?Because it will return error if there is an error, if no error it will return nil in the last line
	book := Book{}

	err := context.BodyParser(&book) //what is BodyParser? It will validate JSON structure
	if err != nil {
		context.Status(http.StatusUnprocessableEntity).JSON( //what is statusunprocessentity? Understand the request but cannot process it
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&book).Error //what using & in book? because &book passed by refrence ||| .Create auto populate it in the database
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not create book"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added"})
	return nil //why return nil? Already answered
}

func (r *Repository) GetBooks(context *fiber.Ctx) error {
	bookModels := &[]models.Books{} //What this &[]models.Book{} means? it will create a slice of book modles from the models package

	err := r.DB.Find(bookModels).Error //The error message means it cant find the book or the connection is error? if it cant find the book it will return empty arr, so the error is more than cant find the book
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get books"})
		return err
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) DeleteBook(context *fiber.Ctx) error {
	bookModel := models.Books{} //why using models.Books() while it dont have data? to specify the type of the model that GORM should delete from the database
	id := context.Params("id")  //what is context.Params? to get the parameter(id data) from the url
	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)
	if err.Error != nil { //why using err.Error while the other function only use err? it need to return result object, so instead of only know the type of error, it will tell you what the error is
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err.Error
	}

	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book deleted successfully",
	})

	return nil
}

func (r *Repository) GetBookByID(context *fiber.Ctx) error {
	id := context.Params("id")
	bookModel := &models.Books{} //why it using &?

	if id == "" {
		context.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}

	fmt.Println("the ID is", id)

	err := r.DB.Where("id=?", id).First(bookModel).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "could not get book"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book fetched successfully",
		"data":    bookModel,
	})

	return nil
}

func main() {
	err := godotenv.Load(".env") //Load the .env file
	if err != nil {              //check error
		log.Fatal(err)
		// log.Fatal("Error loading .env file")
	}

	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"), //what is os.Getenv? it will get the value of the key in the .env file
		Port:     os.Getenv("DB_PORT"), //is DB_HOST is the syntaxt? no, you can change it to anything
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("Could not load the database")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Could not migrate the database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")

}
