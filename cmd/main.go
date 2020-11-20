package main

import (
	"log"
	"net/http"

	"github.com/jrieck1991/golang-di-example/internal/handlers"
)

func main() {

	server := handlers.New()

	log.Fatalln(http.ListenAndServe("localhost:8080", server))
}
