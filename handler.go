package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

//go:embed public/*
var publicFS embed.FS

func publicHandler() http.Handler {
	httpFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatal(err)
	}
	return http.StripPrefix("/public/", http.FileServer(http.FS(httpFS)))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(publicFS, "public/templates/header.tmpl", "public/html/index.html"))
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFS(publicFS, "public/templates/header.tmpl", "public/html/login.html"))
		err := tmpl.ExecuteTemplate(w, "login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")
		fmt.Fprintf(w, fmt.Sprintf("%s %s", email, password))
	}
}
