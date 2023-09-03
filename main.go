package main

import (
	"embed"
	"log"
	"net/http"

	"github.com/matta9001/web-template/db"
	"github.com/matta9001/web-template/handler"
)

//go:embed public/*
var publicFS embed.FS

func main() {

	db.Init()
	handler.InitializeFS(publicFS)

	http.Handle("/public/", handler.PublicHandler())
	http.HandleFunc("/", handler.IndexHandler)

	http.Handle("/profile", handler.AuthMiddleware(http.HandlerFunc(handler.ProfileHandler)))
	http.HandleFunc("/register", handler.RegisterHandler)
	http.HandleFunc("/login", handler.LoginHandler)
	http.HandleFunc("/logout", handler.LogoutHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
