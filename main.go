package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	handler := NewHandler(NewStore())

	log.Printf("Listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, handler.Routes()))
}
