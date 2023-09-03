package handler

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

var publicFS embed.FS

func InitializeFS(fs embed.FS) {
	publicFS = fs
}

func PublicHandler() http.Handler {
	httpFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		log.Fatal(err)
	}
	return http.StripPrefix("/public/", http.FileServer(http.FS(httpFS)))
}
