package main

import (
	"html/template"
	"path/filepath"

	"github.com/Isotop7/liberator/internal/models"
)

type templateData struct {
	Book        *models.Book
	Books       []*models.Book
	LatestBooks []*models.Book
	ActiveBooks []*models.Book
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./assets/templates/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./assets/templates/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./assets/templates/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
