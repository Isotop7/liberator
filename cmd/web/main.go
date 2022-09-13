package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Isotop7/liberator/internal/models"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	_ "github.com/mattn/go-sqlite3"
)

// Liberator struct
type liberator struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	books          *models.BookModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
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

	// Load or create database
	infoLog.Print("Setup database ...")
	db, err := sql.Open("sqlite3", "liberator.db")
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Setup template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Setup session manager
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.IdleTimeout = 60 * time.Minute
	sessionManager.Cookie.Secure = true

	// Setup shared struct
	liberator := &liberator{
		infoLog:        infoLog,
		errorLog:       errorLog,
		books:          &models.BookModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	// TLS config
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
	}

	// Server struct
	srv := &http.Server{
		Addr:         portAddress,
		ErrorLog:     errorLog,
		Handler:      liberator.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  7 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Serve liberator
	errorLog.Fatal(srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"))
}
