package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

//go:embed public/*
var publicFS embed.FS

var store = sessions.NewCookieStore([]byte("secret"))

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

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	email := session.Values["email"]

	tmpl := template.Must(template.ParseFS(publicFS, "public/templates/header.tmpl", "public/html/profile.html"))
	err := tmpl.ExecuteTemplate(w, "profile.html", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tmpl := template.Must(template.ParseFS(publicFS, "public/templates/header.tmpl", "public/html/register.html"))
		err := tmpl.ExecuteTemplate(w, "register.html", nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("password")
		err := register(email, password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		err := login(email, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, _ := store.Get(r, "user-session")
		session.Values["authenticated"] = true
		session.Values["email"] = email
		session.Save(r, w)

		http.Redirect(w, r, "/profile", http.StatusSeeOther)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	session.Options = &sessions.Options{
		MaxAge: -1,
	}
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
