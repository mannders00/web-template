package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/pat"
	"github.com/joho/godotenv"
)

func main() {

	port := ":3000"

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	p := pat.New()

	// OAuth handlers
	SetupProviders()
	p.Get("/auth/{provider}/callback", AuthCallback)
	p.Get("/logout/{provider}", AuthLogout)
	p.Get("/auth/{provider}", AuthHandler)

	// File server
	p.Get("/public/", http.HandlerFunc(http.StripPrefix("/public", http.FileServer(http.Dir("./public"))).ServeHTTP))

	// Template views
	p.Get("/", IndexHandler)

	fmt.Printf("Listening on %s", port)
	err = http.ListenAndServe(port, p)
	if err != nil {
		log.Fatal(err)
	}

}
