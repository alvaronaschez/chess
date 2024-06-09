package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type ChessGame struct {
	whiteReadChannel chan Message
	blackReadChannel chan Message
	whiteWebsocket   *websocket.Conn
	blackWebsocket   *websocket.Conn
}

func NewChessGame(ws *websocket.Conn) *ChessGame {
	game := ChessGame{
		whiteReadChannel: make(chan Message),
		blackReadChannel: make(chan Message),
		whiteWebsocket:   ws,
	}
	return &game
}

func (game *ChessGame) AddWebsocket(ws *websocket.Conn) {
	game.blackWebsocket = ws
}

func (game ChessGame) Start() {
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
