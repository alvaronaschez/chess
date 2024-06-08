package main

import "github.com/gorilla/websocket"

type ChessGame struct {
	WhiteChannel   chan Message
	BlackChannel   chan Message
	WhiteWebsocket *websocket.Conn
	BlackWebsocket *websocket.Conn
}

func NewChessGame(ws *websocket.Conn) *ChessGame {
	game := ChessGame{
		WhiteChannel:   make(chan Message),
		BlackChannel:   make(chan Message),
		WhiteWebsocket: ws,
	}
	return &game
}

func (game *ChessGame) AddWebsocket(ws *websocket.Conn) chan Message {
	game.BlackWebsocket = ws
	return game.BlackChannel
}

func (game ChessGame) Start() {
	go func() {
		turnWhite := true
		game.WhiteWebsocket.WriteJSON(Message{Type: "start", Color: "white"})
		game.BlackWebsocket.WriteJSON(Message{Type: "start", Color: "black"})
		for {
			select {
			case message := <-game.WhiteChannel:
				if turnWhite {
					game.BlackWebsocket.WriteJSON(message)
					turnWhite = false
				}
			case message := <-game.BlackChannel:
				if !turnWhite {
					game.WhiteWebsocket.WriteJSON(message)
					turnWhite = true
				}
			}
		}
	}()
}
