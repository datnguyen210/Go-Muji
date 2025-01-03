package main

import (
	"io/fs"
	"path/filepath"
	"text/template"
	"time"

	"github.com/datnguyen210/go-blog/internal/models"
	"github.com/datnguyen210/go-blog/ui"
)

type templateData struct {
	CurrentYear     int
	Blog            *models.Blog
	Blogs           []*models.Blog
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func formatDate(t time.Time) string {
	if(t.IsZero()){
		return ""
	}
	return t.UTC().Format("2006-01-02 at 15:04")
}

var functions = template.FuncMap{
	"formatDate": formatDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the fs.Glob() to get a slice of all files in the 
	// ui.Files filesystemthat match the pattern "html/pages/*.tmpl".
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
