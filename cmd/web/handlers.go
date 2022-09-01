package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Isotop7/liberator/internal/models"
)

/*// List all Books
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
}*/

func (liberator *liberator) bookCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		liberator.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	id, err := liberator.books.Insert("dnqd", "dojqwdk", "dqwdmkqd", "dqiqdmq", "1234567890", "1234567890123", 555)
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	liberator.infoLog.Printf("Created book with id %v", id)
	http.Redirect(w, r, fmt.Sprintf("/book/view?id=%d", id), http.StatusSeeOther)
}

func (liberator *liberator) bookView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		liberator.notFound(w)
		return
	}

	book, err := liberator.books.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			liberator.notFound(w)
		} else {
			liberator.serverError(w, err)
		}
		return
	}

	data := liberator.newTemplateData(r)
	data.Book = book

	liberator.render(w, http.StatusOK, "book.tmpl", data)
}

func (liberator *liberator) dashboard(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
		return
	} else if path == "/dashboard" {
	} else {
		liberator.notFound(w)
		return
	}

	// Query latest books
	latestBooks, err := liberator.books.Latest(5)
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Get sum page count
	sumPageCount, err := liberator.books.SumPageCount()
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	data := liberator.newTemplateData(r)
	data.LatestBooks = latestBooks
	data.SumPageCount = sumPageCount

	liberator.render(w, http.StatusOK, "dashboard.tmpl", data)
}
