package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/Isotop7/liberator/internal/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Liberator struct
type liberator struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	books         *models.BookModel
	templateCache map[string]*template.Template
}

// Main function
func main() {
	// Setup flags
	port := flag.Int("port", 5000, "Network port")
	flag.Parse()
	portAddress := ":" + strconv.Itoa(*port)

	// Setup loggers
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime)
	infoLog.Printf("Starting liberator on port '%s'", portAddress)

	// Setup template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Setup shared struct
	liberator := &liberator{
		infoLog:       infoLog,
		errorLog:      errorLog,
		books:         &models.BookModel{DB: nil},
		templateCache: templateCache,
	}

	// Load or create database
	infoLog.Print("Setup database ...")
	db, err := gorm.Open("sqlite3", "liberator.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	db.AutoMigrate(&models.Book{})
	liberator.books.DB = db
	defer db.Close()

	// Server struct
	srv := &http.Server{
		Addr:     portAddress,
		ErrorLog: errorLog,
		Handler:  liberator.routes(),
	}

	// Serve liberator
	errorLog.Fatal(srv.ListenAndServe())
}
