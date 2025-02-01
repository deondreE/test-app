package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var ctx = context.Background()
var clients = make(map[*websocket.Conn]bool)
var clientsMutex = sync.Mutex{}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "admin",
		DB:       0,
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		clientsMutex.Lock()
		clients[conn] = true
		clientsMutex.Unlock()

		defer func() {
			clientsMutex.Lock()
			delete(clients, conn)
			clientsMutex.Unlock()
		}()

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			clientsMutex.Lock()
			for client := range clients {
				if err := client.WriteMessage(msgType, msg); err != nil {
					fmt.Println("Error sending message to client:", err)
					client.Close()
					delete(clients, client) // remove the client if there's an error
				}
			}
			clientsMutex.Unlock()
		}
	})

	http.HandleFunc("/decode", func(w http.ResponseWriter, r *http.Request) {
		valueJSON, err := rdb.Get(ctx, "Peter").Result()
		if err != nil {
			http.Error(w, "Failed to retrieve data from Redis", http.StatusInternalServerError)
			return
		}

		var user User
		err = json.Unmarshal([]byte(valueJSON), &user)
		if err != nil {
			http.Error(w, "Failed to deserialize user data", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "User: %s %s, Age: %d", user.Firstname, user.Lastname, user.Age)
	})

	http.HandleFunc("/encode", func(w http.ResponseWriter, r *http.Request) {
		peter := User{
			Firstname: "John",
			Lastname:  "Doe",
			Age:       25,
		}

		peterJSON, err := json.Marshal(peter)
		if err != nil {
			http.Error(w, "Failed to serialize user data", http.StatusInternalServerError)
			return
		}

		err = rdb.Set(ctx, "Peter", peterJSON, 0).Err()
		if err != nil {
			http.Error(w, "Failed to save data to Redis", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Peter data stored in Redis successfully!")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "test.html")
	})

	http.ListenAndServe(":8080", nil)
}
