package main

import (
	"log"
	"net/http"
)

func main() {

	initDB()

	http.Handle("/public/", publicHandler())
	http.HandleFunc("/", indexHandler)

	http.Handle("/profile", authMiddleware(http.HandlerFunc(profileHandler)))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
