package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"text/template"
)

//go:embed public/*
var publicFS embed.FS

func main() {

	http.Handle("/public/", publicHandler())
	http.HandleFunc("/", indexHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}

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
