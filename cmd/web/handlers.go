package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Isotop7/liberator/internal/models"
	"github.com/julienschmidt/httprouter"
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

func (liberator *liberator) dashboard(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
		return
	}

	// Query latest books
	latestBooks, err := liberator.books.Latest(5)
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Get sum page count
	sumPageCount := liberator.books.SumPageCount()

	data := liberator.newTemplateData(r)
	data.LatestBooks = latestBooks
	data.SumPageCount = sumPageCount

	liberator.render(w, http.StatusOK, "dashboard.tmpl", data)
}

func (liberator *liberator) bookCreate(w http.ResponseWriter, r *http.Request) {
	data := liberator.newTemplateData(r)

	liberator.render(w, http.StatusOK, "bookCreate.tmpl", data)
}

func (liberator *liberator) bookCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	title := r.PostForm.Get("title")
	author := r.PostForm.Get("author")
	language := r.PostForm.Get("language")
	category := r.PostForm.Get("category")
	isbn10 := r.PostForm.Get("isbn10")
	isbn13 := r.PostForm.Get("isbn13")
	review := r.PostForm.Get("review")

	pagecount, err := strconv.Atoi(r.PostForm.Get("pagecount"))
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}
	// Check for invalid page count
	if pagecount < 1 {
		liberator.clientError(w, http.StatusBadRequest)
	}

	rating, err := strconv.Atoi(r.PostForm.Get("rating"))
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}
	// Check for invalid rating
	if rating < 1 || rating > 10 {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Insert element to database
	id, err := liberator.books.Insert(title, author, language, category, isbn10, isbn13, pagecount, rating, review)
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Redirect to view
	http.Redirect(w, r, fmt.Sprintf("/book/view/%d", id), http.StatusSeeOther)
}

func (liberator *liberator) bookView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

	liberator.render(w, http.StatusOK, "bookView.tmpl", data)
}
