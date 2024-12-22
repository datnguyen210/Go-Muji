package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/datnguyen210/go-blog/internal/models"
)

type templateData struct {
	CurrentYear int
	Blog  *models.Blog
	Blogs []*models.Blog
	Form any
	Flash string
}

func formatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04")
}

var functions = template.FuncMap{
	"formatDate": formatDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob() function to get a slice of all file that
	// match the pattern "./ui/html/pages/*.tmpl".
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
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
