package main

import (
	"log"
	"net/http"
)

func main() {
	setupAPI()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func setupAPI() {
	manager := NewManager()
	http.HandleFunc("/ws", manager.serveWS)
}
