package main

import (
	"log"
	"net/http"
)

func main() {

	http.Handle("/public/", publicHandler())
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
