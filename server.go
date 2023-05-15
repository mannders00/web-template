package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.PathPrefix("/public/").HandlerFunc(getPublic)

	fmt.Println("starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("public/templates/header.tmpl", "public/html/index.html"))
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPublic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	filePath := filepath.Join(".", path)

	if strings.HasPrefix(path, "/public/css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasPrefix(path, "/public/js") {
		w.Header().Set("Content-Type", "application/javascript")
	} else if strings.HasPrefix(path, "/public/images") {
		w.Header().Set("Content-Type", "image/jpeg")
	} else if strings.HasPrefix(path, "/data/images") {
		w.Header().Set("Content-Type", "image/png")
	} else if strings.HasPrefix(path, "/public/fonts") {
		w.Header().Set("Content-Type", "fonts/font")
	}

	http.ServeFile(w, r, filePath)
}
