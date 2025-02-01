package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println(AllRooms.Map)
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Fatal(err)
					client.Conn.Close()
				}

			}
		}
	}
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]
	if !ok {
		log.Println("RoomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error:", err)
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		log.Println(msg.Message)

		broadcast <- msg

	}

}

// Send messages
type Message struct {
	ID       string `json:"id"`
	Msg      string `json:"msg"`
	SentTime string `json:"sent_time"`
	Session  string `json:"session"`
	Username string `json:"username"`
}

var ctx = context.Background()
var clients = make(map[*websocket.Conn]string) // map clients to their session IDs
var sessionStore = make(map[string]string)     // map session IDs to user IDs or valid sessions
var clientsMutex = sync.Mutex{}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:1420")
}

func GenerateUniqueId() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func GenerateSessionId() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}

func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Upgrade HTTP request to WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade WebSocket", http.StatusInternalServerError)
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	sessionId := GenerateSessionId()
	// TODO: get username from user
	clientsMutex.Lock()
	clients[conn] = sessionId
	clientsMutex.Unlock()
	sessionStore[sessionId] = "user"

	for {
		msgType, msg, err := conn.ReadMessage()
		id := GenerateUniqueId()
		if err != nil {
			return
		}

		var rMessage struct {
			Username string `json:"username"`
			Msg      string `json:"message"`
		}

		if err := json.Unmarshal(msg, &rMessage); err != nil {
			fmt.Println("Error unmarshalling message:", err)
			continue
		}

		message := Message{ID: id, Username: rMessage.Username, Msg: rMessage.Msg, Session: sessionId, SentTime: time.Now().Format("2017-09-07 17:06:06")}
		fmt.Printf("Received message from %s at %s\n", message.Msg, message.SentTime)

		resJson, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error marshaling message:", err)
			continue
		}

		clientsMutex.Lock()
		for client := range clients {
			if err := client.WriteMessage(msgType, resJson); err != nil {
				fmt.Println("Error sending message to client:", err)
				client.Close()
				delete(clients, client) // remove the client if there's an error
			}
		}
		clientsMutex.Unlock()
	}
}
