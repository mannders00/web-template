package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
	ory "github.com/ory/client-go"
)

//go:embed public/*
var publicFS embed.FS

type App struct {
	ory *ory.APIClient
	db  *sql.DB
}

func main() {

	// DB Initialization
	newDB, err := sql.Open("sqlite3", "./db.db")
	if err != nil {
		log.Fatal(err)
	}
	defer newDB.Close()

	statement, err := newDB.Prepare("CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, active INTEGER, streak TEXT)")
	_, err = statement.Exec()
	if err != nil {
		log.Fatal(err)
	}

	// Ory Authentication
	proxyPort := "4000"
	c := ory.NewConfiguration()
	c.Servers = ory.ServerConfigurations{{URL: fmt.Sprintf("http://localhost:%s/.ory", proxyPort)}}
	app := &App{
		ory: ory.NewAPIClient(c),
		db:  newDB,
	}

	// HTTP Server
	mux := http.NewServeMux()
	mux.Handle("/", app.sessionMiddleware(app.getIndexHandler()))
	mux.Handle("/public/", publicHandler())
	mux.HandleFunc("/clicked", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postClickHandler(w, r)
		}
	})
	mux.HandleFunc("/test", getTestHandler)

	err = http.ListenAndServe(":3000", mux)
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

func postClickHandler(w http.ResponseWriter, r *http.Request) {
	html := `
		<button hx-post="/clicked" hx-swap="outerHTML" class="btn btn-primary">fuck you</button>
	`
	_, err := fmt.Fprintf(w, html)
	if err != nil {
		log.Fatal(err)
	}
}

func getTestHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFS(publicFS, "public/templates/header.tmpl", "public/html/index.html"))
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (app *App) getIndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFS(publicFS, "public/templates/header.tmpl", "public/html/index.html"))
		session := getSession(r.Context())
		fmt.Println(session.Id)
		sessionJSON, err := json.Marshal(session)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "index.html", string(sessionJSON))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
