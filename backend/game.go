package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ChessGame struct {
	whiteWebsocket   *websocket.Conn
	blackWebsocket   *websocket.Conn
	whiteReadChannel chan Message
	blackReadChannel chan Message
}

func NewChessGame(ws *websocket.Conn) *ChessGame {
	game := ChessGame{
		whiteWebsocket:   ws,
		whiteReadChannel: make(chan Message),
		blackReadChannel: make(chan Message),
	}
	return &game
}

func (game *ChessGame) Join(ws *websocket.Conn) {
	game.blackWebsocket = ws
	go forwardFromWebsocketToChannel(game.whiteWebsocket, game.whiteReadChannel)
	go forwardFromWebsocketToChannel(game.blackWebsocket, game.blackReadChannel)
	go func() {
		turnWhite := true
		game.whiteWebsocket.WriteJSON(Message{Type: "start", Color: "white"})
		game.blackWebsocket.WriteJSON(Message{Type: "start", Color: "black"})
		for {
			select {
			case message := <-game.whiteReadChannel:
				if turnWhite {
					game.blackWebsocket.WriteJSON(message)
					turnWhite = false
				}
			case message := <-game.blackReadChannel:
				if !turnWhite {
					game.whiteWebsocket.WriteJSON(message)
					turnWhite = true
				}
			}
		}
	}()
}

func forwardFromWebsocketToChannel(ws *websocket.Conn, ch chan Message) {
	defer ws.Close()
	for {
		message := Message{}
		err := ws.ReadJSON(&message)

		if err != nil {
			log.Println(err)
			return
		}

		ch <- message
	}
}
