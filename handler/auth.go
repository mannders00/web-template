package handler

import (
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/matta9001/web-template/db"
)

var store = sessions.NewCookieStore([]byte("secret"))

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "user-session")
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	email := session.Values["email"]

	tmpl := template.Must(template.ParseFS(publicFS, "public/templates/header.tmpl", "public/html/profile.html"))
	err := tmpl.ExecuteTemplate(w, "profile.html", email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
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
		err := db.Register(email, password)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		err := db.Login(email, password)
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	session.Options = &sessions.Options{
		MaxAge: -1,
	}
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
