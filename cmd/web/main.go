package main

import (
	"log"
	"net/http"

	"github.com/iradukunda1/web-socket/internals/handlers"
)

func main() {
	server := routes()
	log.Println("starting websocket listener channel")

	go handlers.ListenToWsChannel()

	log.Print("Listening on :8000")
	http.ListenAndServe(":8000", server)
}
