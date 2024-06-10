package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ChessGame struct {
	whiteWebsocket   *websocket.Conn
	blackWebsocket   *websocket.Conn
	whiteReadChannel chan Message
	blackReadChannel chan Message
	done             chan struct{}
	finished         bool
	mu               sync.Mutex
}

type Message struct {
	Type      string `json:"type" validate:"required,oneof=start move"`
	Color     string `json:"color" validate:"oneof=white black,required_if=Type start"`
	From      string `json:"from" validate:"required_if=Type move"`
	To        string `json:"to" validate:"required_if=Type move"`
	Promotion string `json:"promotion" validate:"oneof=q r b k,required_if=Type move"`
}

func NewChessGame(ws *websocket.Conn) *ChessGame {
	game := ChessGame{
		whiteWebsocket:   ws,
		whiteReadChannel: make(chan Message),
		blackReadChannel: make(chan Message),
		done:             make(chan struct{}),
		finished:         false,
	}
	return &game
}

func (game *ChessGame) Join(ws *websocket.Conn) {
	game.blackWebsocket = ws
	go game.forwardFromWebsocketToChannel("white")
	go game.forwardFromWebsocketToChannel("black")
	go func() {
		turnWhite := true
		game.whiteWebsocket.WriteJSON(Message{Type: "start", Color: "white"})
		game.blackWebsocket.WriteJSON(Message{Type: "start", Color: "black"})
		for {
			select {
			case <-game.done:
				log.Println("game finished")
				return
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

func (game *ChessGame) forwardFromWebsocketToChannel(color string) {
	var ws *websocket.Conn
	var ch chan Message
	if color == "white" {
		ws = game.whiteWebsocket
		ch = game.whiteReadChannel
	} else {
		ws = game.blackWebsocket
		ch = game.blackReadChannel
	}
	defer ws.Close()
	for {
		message := Message{}
		err := ws.ReadJSON(&message)

		if err != nil {
			log.Println(err)
			// lock before reading game.finished
			game.mu.Lock()
			defer game.mu.Unlock()
			if !game.finished {
				game.finished = true
				close(game.done)
			}
			return
		}

		ch <- message
	}
}
