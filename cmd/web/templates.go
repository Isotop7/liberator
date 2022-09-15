package main

import (
	"database/sql"
	"html/template"
	"io/fs"
	"path/filepath"

	"github.com/Isotop7/liberator/internal/models"
	"github.com/Isotop7/liberator/ui"
)

type templateData struct {
	CurrentYear     int
	Book            *models.Book
	Books           []*models.Book
	LatestBooks     []*models.Book
	ActiveBooks     []*models.Book
	BookIsAssigned  bool
	SumPageCount    int
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("02.01.2006 15:04")
	} else {
		return ""
	}
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
