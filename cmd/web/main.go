package main

import (
	"flag"
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
	errorLog *log.Logger
	infoLog  *log.Logger
	books    *models.BookModel
}

// Main function
func main() {
	// Setup flags
	port := flag.Int("port", 5000, "Network port")
	flag.Parse()
	portAddress := ":" + strconv.Itoa(*port)

	// Setup loggers
	infoLog := log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "[ERROR] ", log.Ldate|log.Ltime)

	// Setup shared struct
	liberator := &liberator{
		infoLog:  infoLog,
		errorLog: errorLog,
		books:    &models.BookModel{DB: nil},
	}

	infoLog.Printf("Starting liberator on port '%s'", portAddress)
	infoLog.Print("Setup database ...")

	db, err := gorm.Open("sqlite3", "liberator.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	liberator.books.DB = db
	defer db.Close()

	// Server struct
	lbrtr := &http.Server{
		Addr:     portAddress,
		ErrorLog: errorLog,
		Handler:  liberator.routes(),
	}

	// Serve liberator
	errorLog.Fatal(lbrtr.ListenAndServe())
}
