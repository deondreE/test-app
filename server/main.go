package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/deondreE/test-app/server"
)

func main() {
	server.AllRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)
	http.HandleFunc("/message", server.SendMessageHandler)

	log.Println("starting Server on  Port 8080")
	fmt.Println(" ")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal((err))
	}
}
