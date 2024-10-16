package main

import (
	"path/filepath"
	"text/template"

	"github.com/datnguyen210/go-blog/internal/models"
)

type templateData struct {
	Blog  *models.Blog
	Blogs []*models.Blog
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
	// Loop through the page filepaths one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'home.tmpl')
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Create a slice containing the filepaths for 
		// our base templates as well as the current page.
		files := []string{
			"./ui/html/base.tmpl",
			"./ui/html/partials/nav.tmpl",
			page,
		}
		// Parse the files into a template set.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
