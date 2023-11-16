package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/markbates/goth"

	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/apple"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/twitterv2"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("public/templates/header.tmpl", "public/html/index.html"))
	err := tmpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// OAuth2 Support

func SetupProviders() {
	goth.UseProviders(
		twitterv2.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), "http://localhost:3000/auth/twitterv2/callback"),
		apple.New(os.Getenv("APPLE_KEY"), os.Getenv("APPLE_SECRET"), "http://localhost:3000/auth/apple/callback", nil, apple.ScopeName, apple.ScopeEmail),
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:3000/auth/google/callback"),
	)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	if gothUser, err := gothic.CompleteUserAuth(w, r); err == nil {
		tmpl := template.Must(template.ParseFiles("public/templates/header.tmpl", "public/html/index.html"))
		tmpl.ExecuteTemplate(w, "index.html", gothUser)
	} else {
		gothic.BeginAuthHandler(w, r)
	}
}

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	tmpl := template.Must(template.ParseFiles("public/templates/header.tmpl", "public/html/index.html"))
	tmpl.ExecuteTemplate(w, "index.html", user)
}

func AuthLogout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}
