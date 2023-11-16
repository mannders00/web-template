package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	port := ":8080"

	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))

	http.HandleFunc("/", IndexHandler)

	fmt.Printf("Listening on %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
