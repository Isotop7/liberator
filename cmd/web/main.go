package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DB *gorm.DB

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

	DB = db
}

// Main function
func main() {
	log.Println("Starting liberator ...")

	log.Println("Setup database ...")
	connectDatabase()

	// Setup mux
	mux := http.NewServeMux()

	// Serve static files
	fileServer := http.FileServer(http.Dir("./assets/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Index handler
	mux.HandleFunc("/", dashboard)
	mux.HandleFunc("/dashboard", dashboard)

	// Books
	mux.HandleFunc("/book", book)

	// Serve liberator
	log.Fatal(http.ListenAndServe(":5000", mux))
}
