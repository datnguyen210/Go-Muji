package main

import (
	"html/template"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	
	ts, err := template.ParseFiles(files...) // destructure
	if err != nil {
		log.Println(err.Error())
		// return a general error message instead of letting 
		// the users know more info to avoid attacker
		http.Error(w, "Internal Server Error", 500) 
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func blogView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific blog..."))
}

func blogCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new blog..."))
}