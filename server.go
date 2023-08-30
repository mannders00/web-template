package main

import (
	"log"
	"net/http"
)

func main() {

	initDB()

	http.Handle("/public/", publicHandler())
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/register", registerHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
