package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

// Shelve
type Shelve struct {
	gorm.Model
	Location string `json:"location"`
	Content  []Book `json:"content"`
}

// Book
type Book struct {
	gorm.Model
	Title     string `json:"title"`
	Author    string `json:"author"`
	Language  string `json:"language"`
	Category  string `json:"category"`
	ISBN10    string `json:"isbn10" binding:"len=10"`
	ISBN13    string `json:"isbn13" binding:"len=13"`
	PageCount int    `json:"page_count"`
	Rating    int    `json:"rating"`
}

// Connect to database
func connectDatabase() {
	db, err := gorm.Open("sqlite3", "liberator.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&Book{})
	db.AutoMigrate(&Shelve{})

	DB = db
}

// List all Books
func listBooksEndpoint(ctx *gin.Context) {
	var books = []Book{}
	err := DB.Find(&books)
	if err != nil {
		ctx.JSON(http.StatusOK, books)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "error getting list of books",
		})
	}
}

// List single book
func listBookEndpoint(ctx *gin.Context) {
	// Get parameter from request
	idParam := ctx.Param("id")

	// Parse id to int
	id, err := strconv.Atoi(idParam)
	// If no int
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}
	// If negative value
	if id <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}

	// Create and find book
	var book = Book{}
	dbErr := DB.Where("id = ?", id).First(&book)

	// If book was found
	if dbErr != nil && book.ID > 0 {
		ctx.JSON(http.StatusOK, book)
		return
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}
}

// Create single book
func createBookEndpoint(ctx *gin.Context) {
	var newElement Book

	// Bind body data to element
	err := ctx.BindJSON(&newElement)
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
		}
		return
	}

	// Check if element with id is already present in db and fail
	var duplicateBook Book
	_ = DB.Where("id = ?", newElement.ID).First(&duplicateBook)
	if duplicateBook.ID == newElement.ID {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "duplicate element with id",
			"data":    newElement,
		})
		return
	} else {
		// Create element in db
		DB.Create(&newElement)
		ctx.JSON(http.StatusCreated, newElement)
	}
}

// Update single book
func updateBookEndpoint(ctx *gin.Context) {
	// Get parameter from request
	idParam := ctx.Param("id")

	// Parse id to int
	id, err := strconv.Atoi(idParam)
	// If no int
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}
	// If negative value
	if id <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}

	// Create and find book
	var book = Book{}
	dbErr := DB.Where("id = ?", id).First(&book)

	// If book was found
	if dbErr != nil && book.ID > 0 {
		bindErr := ctx.BindJSON(&book)
		if bindErr != nil {
			ctx.JSON(http.StatusMethodNotAllowed, gin.H{
				"message": bindErr.Error(),
			})
			return
		}
		DB.Save(&book)
		ctx.JSON(http.StatusOK, book)
		return
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}
}

func deleteBookEndpoint(ctx *gin.Context) {
	// Get parameter from request
	idParam := ctx.Param("id")

	// Parse id to int
	id, err := strconv.Atoi(idParam)
	// If no int
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}
	// If negative value
	if id <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "book not found",
		})
		return
	}

	// Create and find book
	var book = Book{}
	dbErr := DB.Where("id = ?", id).Delete(&book)
	if dbErr != nil {
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "book with id '" + idParam + "' deleted",
		})
	} else {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": dbErr,
		})
	}
}

// List all shelves
func listShelvesEndpoint(ctx *gin.Context) {
	var shelves = []Shelve{}
	err := DB.Find(&shelves)
	if err != nil {
		ctx.JSON(http.StatusOK, shelves)
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "error getting list of shelves",
		})
	}
}

// List single shelve
func listShelveEndpoint(ctx *gin.Context) {
	// Get parameter from request
	idParam := ctx.Param("id")

	// Parse id to int
	id, err := strconv.Atoi(idParam)
	// If no int
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "shelve not found",
		})
		return
	}
	// If negative value
	if id <= 0 {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "shelve not found",
		})
		return
	}

	// Create and find shelve
	var shelve = Shelve{}
	dbErr := DB.Where("id = ?", id).First(&shelve)

	// If shelve was found
	if dbErr != nil && shelve.ID > 0 {
		ctx.JSON(http.StatusOK, shelve)
		return
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "shelve not found",
		})
		return
	}
}

// Create single shelve
// TODO: Handle books in content body
func createShelveEndpoint(ctx *gin.Context) {
	var newElement Shelve

	// Bind body data to element
	err := ctx.BindJSON(&newElement)
	if err != nil {
		switch err.(type) {
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid request",
			})
		}
		return
	}

	// Check if element with id is already present in db and fail
	var duplicateShelve Shelve
	_ = DB.Where("id = ?", newElement.ID).First(&duplicateShelve)
	if duplicateShelve.ID == newElement.ID {
		ctx.JSON(http.StatusConflict, gin.H{
			"message": "duplicate element with id",
			"data":    newElement,
		})
		return
	} else {
		// Create element in db
		DB.Create(&newElement)
		ctx.JSON(http.StatusCreated, newElement)
	}
}

// Main function
func main() {
	log.Println("Starting liberator ...")

	log.Println("Setup database ...")
	connectDatabase()

	// Setup handlers
	router := gin.Default()

	// Shelves
	router.GET("/shelves", listShelvesEndpoint)
	router.GET("/shelves/:id", listShelveEndpoint)
	router.POST("/shelves", createShelveEndpoint)

	// Books
	router.GET("/books", listBooksEndpoint)
	router.POST("/books", createBookEndpoint)
	router.GET("/books/:id", listBookEndpoint)
	router.PUT("/books/:id", updateBookEndpoint)
	router.DELETE("/books/:id", deleteBookEndpoint)

	// Serve API
	log.Fatal(router.Run(":5000"))
}
