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
	Review    string
	validator.Validator
}

type userSignupForm struct {
	Name     string
	Email    string
	Password string
	validator.Validator
}

type userLoginForm struct {
	Email    string
	Password string
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

	// Get user id
	id := liberator.sessionManager.GetInt(r.Context(), "authenticatedUserID")
	// Get sum page count
	sumPageCount := liberator.booksUsersAssignment.SumPageCount(id)
	// Get active books
	activeBooks, err := liberator.booksUsersAssignment.GetActiveBooks(id)
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	data := liberator.newTemplateData(r)
	data.LatestBooks = latestBooks
	data.ActiveBooks = activeBooks
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

	// Init struct form
	form := bookCreateForm{
		Title:     r.PostForm.Get("title"),
		Author:    r.PostForm.Get("author"),
		Language:  r.PostForm.Get("language"),
		Category:  r.PostForm.Get("category"),
		ISBN10:    r.PostForm.Get("isbn10"),
		ISBN13:    r.PostForm.Get("isbn13"),
		Pagecount: pagecount,
		Review:    r.PostForm.Get("review"),
	}

	// Validate title
	form.CheckField(validator.NotBlank(form.Title), "title", validator.ValueMustNotBeEmpty)
	form.CheckField(validator.MaxChars(form.Title, 255), "title", validator.ValueMustNotBeLongerThan(255))

	// Validate author
	form.CheckField(validator.NotBlank(form.Author), "author", validator.ValueMustNotBeEmpty)
	form.CheckField(validator.MaxChars(form.Author, 255), "author", validator.ValueMustNotBeLongerThan(255))

	// Validate category
	form.CheckField(validator.NotBlank(form.Category), "category", validator.ValueMustNotBeEmpty)
	form.CheckField(validator.MaxChars(form.Category, 255), "category", validator.ValueMustNotBeLongerThan(255))

	// Check for invalid page count
	form.CheckField(validator.GreaterThan(form.Pagecount, 0), "pagecount", validator.ValueMustBeGreaterThan(0))

	// Check for invalid ISBN-10
	form.CheckField(validator.IsValidISBN10(form.ISBN10), "isbn10", validator.InvalidISBN)

	// Check for invalid ISBN-13
	form.CheckField(validator.IsValidISBN13(form.ISBN13), "isbn13", validator.InvalidISBN)

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

func (liberator *liberator) bookAssignPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Convert bookID
	bookID, err := strconv.Atoi(r.PostForm.Get("bookID"))
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	book, err := liberator.books.Get(bookID)

	// If error while getting book
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Check if book is already assigned
	alreadyAssigned := liberator.booksUsersAssignment.IsCurrentlyAssigned(int(book.ID))

	if alreadyAssigned {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Get user id
	userid := liberator.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	// Insert assignment into database
	_, err = liberator.booksUsersAssignment.Assign(bookID, userid)

	// If error while inserting
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Save success message to session data
	liberator.sessionManager.Put(r.Context(), "flash", "Buch wurde ausgeliehen!")

	// Redirect to view
	http.Redirect(w, r, fmt.Sprintf("/book/view/%d", bookID), http.StatusSeeOther)
}

func (liberator *liberator) bookUnassignPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Convert bookID
	bookID, err := strconv.Atoi(r.PostForm.Get("bookID"))
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	book, err := liberator.books.Get(bookID)

	// If error while getting book
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Get user id
	userid := liberator.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	// Check if book is assigned to user
	isAssignedToUser := liberator.booksUsersAssignment.IsCurrentlyAssignedToUser(int(book.ID), userid)

	if !isAssignedToUser {
		liberator.clientError(w, http.StatusBadRequest)
	}

	// Insert assignment into database
	_, err = liberator.booksUsersAssignment.UpdateAssignmentState(bookID, userid, models.Inactive)

	// If error while inserting
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	// Save success message to session data
	liberator.sessionManager.Put(r.Context(), "flash", "Buch wurde zurückgegeben!")

	// Redirect to view
	http.Redirect(w, r, fmt.Sprintf("/book/view/%d", bookID), http.StatusSeeOther)
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

	// Check book assignment state
	bookIsAssigned := liberator.booksUsersAssignment.IsCurrentlyAssigned(id)

	data := liberator.newTemplateData(r)
	data.Book = book
	data.BookIsAssigned = bookIsAssigned

	liberator.render(w, http.StatusOK, "bookView.tmpl", data)
}

func (liberator *liberator) userSignup(w http.ResponseWriter, r *http.Request) {
	data := liberator.newTemplateData(r)
	data.Form = userSignupForm{}
	liberator.render(w, http.StatusOK, "signup.tmpl", data)
}

func (liberator *liberator) userSignupPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	form := userSignupForm{
		Name:     r.PostForm.Get("name"),
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Name), "name", validator.ValueMustNotBeEmpty)
	form.CheckField(validator.NotBlank(form.Email), "email", validator.ValueMustNotBeEmpty)
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", validator.ValueInvalidEmail)
	form.CheckField(validator.NotBlank(form.Password), "password", validator.ValueMustNotBeEmpty)
	form.CheckField(validator.MinChars(form.Password, 8), "password", validator.ValueMustBeLongerThan(8))

	if !form.Valid() {
		data := liberator.newTemplateData(r)
		data.Form = form
		liberator.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	err = liberator.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Die Email-Adresse wird bereits verwendet")

			data := liberator.newTemplateData(r)
			data.Form = form
			liberator.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			liberator.serverError(w, err)
		}
		return
	}

	liberator.sessionManager.Put(r.Context(), "flash", "Die Registrierung war erfolgreich. Du kannst dich nun anmelden.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (liberator *liberator) userLogin(w http.ResponseWriter, r *http.Request) {
	data := liberator.newTemplateData(r)
	data.Form = userLoginForm{}
	liberator.render(w, http.StatusOK, "login.tmpl", data)
}

func (liberator *liberator) userLoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		liberator.clientError(w, http.StatusBadRequest)
	}

	form := userLoginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(validator.NotBlank(form.Email), "email", validator.ValueMustNotBeEmpty)
	form.CheckField(validator.Matches(form.Email, validator.EmailRegex), "email", validator.ValueInvalidEmail)
	form.CheckField(validator.NotBlank(form.Password), "password", validator.ValueMustNotBeEmpty)

	if !form.Valid() {
		data := liberator.newTemplateData(r)
		data.Form = form
		liberator.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	id, err := liberator.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email oder Passwort nicht korrekt")

			data := liberator.newTemplateData(r)
			data.Form = form
			liberator.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			liberator.serverError(w, err)
		}
		return
	}

	err = liberator.sessionManager.RenewToken(r.Context())
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	liberator.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (liberator *liberator) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := liberator.sessionManager.RenewToken(r.Context())
	if err != nil {
		liberator.serverError(w, err)
		return
	}

	liberator.sessionManager.Remove(r.Context(), "authenticatedUserID")
	liberator.sessionManager.Put(r.Context(), "flash", "Du wurdest abgemeldet!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
