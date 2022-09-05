package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Isotop7/liberator/internal/models"
	"github.com/Isotop7/liberator/internal/validator"
	"github.com/julienschmidt/httprouter"
)

type bookCreateForm struct {
	Title     string
	Author    string
	Language  string
	Category  string
	ISBN10    string
	ISBN13    string
	Pagecount int
	Rating    int
	Review    string
	validator.Validator
}

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

func (liberator *liberator) searchView(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	query := r.PostForm.Get("query")

	if query == "" {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	books, err := liberator.books.Search(query)
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	data := liberator.newTemplateData(r)
	data.Books = books

	liberator.render(w, http.StatusOK, "search.tmpl", data)
}

func (liberator *liberator) bookCreate(w http.ResponseWriter, r *http.Request) {
	data := liberator.newTemplateData(r)

	data.Form = bookCreateForm{}

	liberator.render(w, http.StatusOK, "bookCreate.tmpl", data)
}

func (liberator *liberator) bookCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Convert pagecount
	pagecount, err := strconv.Atoi(r.PostForm.Get("pagecount"))
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Convert rating
	rating, err := strconv.Atoi(r.PostForm.Get("rating"))
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Init struct form
	form := bookCreateForm{
		Title:     r.PostForm.Get("title"),
		Author:    r.PostForm.Get("author"),
		Language:  r.PostForm.Get("language"),
		Category:  r.PostForm.Get("category"),
		ISBN10:    r.PostForm.Get("isbn10"), //TODO: Add format check
		ISBN13:    r.PostForm.Get("isbn13"), //TODO: Add checksum method
		Pagecount: pagecount,
		Rating:    rating,
		Review:    r.PostForm.Get("review"),
	}

	// Validate title
	form.CheckValue(validator.NotBlank(form.Title), "title", validator.ValueMustNotBeEmpty)
	form.CheckValue(validator.MaxChars(form.Title, 255), "title", validator.ValueMustNotBeLongerThan(255))

	// Validate author
	form.CheckValue(validator.NotBlank(form.Author), "author", validator.ValueMustNotBeEmpty)
	form.CheckValue(validator.MaxChars(form.Author, 255), "author", validator.ValueMustNotBeLongerThan(255))

	// Validate category
	form.CheckValue(validator.NotBlank(form.Category), "category", validator.ValueMustNotBeEmpty)
	form.CheckValue(validator.MaxChars(form.Category, 255), "category", validator.ValueMustNotBeLongerThan(255))

	// Check for invalid page count
	form.CheckValue(validator.GreaterThan(form.Pagecount, 0), "pagecount", validator.ValueMustBeGreaterThan(0))

	// Check for invalid rating
	form.CheckValue(validator.InBounds(form.Rating, 1, 10), "rating", validator.ValueMustBeInRange(1, 10))

	// If errors are found, show them and redirect to form
	if !form.Valid() {
		data := liberator.newTemplateData(r)
		data.Form = form
		liberator.render(w, http.StatusUnprocessableEntity, "bookCreate.tmpl", data)
		return
	}

	// Insert element to database
	id, err := liberator.books.Insert(
		form.Title,
		form.Author,
		form.Language,
		form.Category,
		form.ISBN10,
		form.ISBN13,
		form.Pagecount,
		form.Rating,
		form.Review)

	// If error while inserting
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Save success message to session data
	liberator.sessionManager.Put(r.Context(), "flash", "Buch wurde erfolgreich erstellt!")

	// Redirect to view
	http.Redirect(w, r, fmt.Sprintf("/book/view/%d", id), http.StatusSeeOther)
}

func (liberator *liberator) bookView(w http.ResponseWriter, r *http.Request) {
	// Parse parameters
	params := httprouter.ParamsFromContext(r.Context())

	// Convert ID from request
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		liberator.notFound(w)
		return
	}

	// Query book by id
	book, err := liberator.books.Get(id)
	// Check for errors
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

func (liberator *liberator) userSignup(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "HTML form for signup")
}

func (liberator *liberator) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "POST handler for signup")
}

func (liberator *liberator) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "HTML login form")
}

func (liberator *liberator) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "POST handler for login")
}

func (liberator *liberator) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "POST handler for logout")
}
