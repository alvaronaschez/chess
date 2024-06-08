package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

var game *ChessGame

func ping(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	defer ws.Close()

	var ch chan Message
	if game == nil {
		game = NewChessGame(ws)
		ch = game.WhiteChannel
	} else {
		game.AddWebsocket(ws)
		ch = game.BlackChannel
		game.Start()
		game = nil
	}

	for {
		message := Message{}
		err = ws.ReadJSON(&message)

		if err != nil {
			log.Println(err)
			return
		}

		ch <- message

		log.Printf("Message: %s", message)
	}
}

func main() {
	fmt.Println("Listening at port 5555")
	http.HandleFunc("/ping", ping)
	log.Fatal(http.ListenAndServe(":5555", nil))
}
