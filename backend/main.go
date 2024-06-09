package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

var game *ChessGame

func ping(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	if game == nil {
		game = NewChessGame(ws)
	} else {
		game.AddWebsocket(ws)
		game.Start()
		game = nil
	}
}

func main() {
	fmt.Println("Listening at port 5555")
	http.HandleFunc("/ping", ping)
	log.Fatal(http.ListenAndServe(":5555", nil))
}
