package main

import (
	"errors"

	"github.com/gorilla/websocket"
)

type ChessGame struct {
	whiteWebsocket *websocket.Conn
	blackWebsocket *websocket.Conn
}

type Message struct {
	Type      string `json:"type" validate:"required,oneof=start move error"`
	Color     string `json:"color" validate:"oneof=white black,required_if=Type start|requited_if=Type move"`
	From      string `json:"from" validate:"required_if=Type move"`
	To        string `json:"to" validate:"required_if=Type move"`
	Promotion string `json:"promotion" validate:"oneof=q r b k,required_if=Type move"`
	WhiteTime int    `json:"whiteTime" validate:"required_if=Type move|required_if=Type start"`
	BlackTime int    `json:"blackTime" validate:"required_if=Type move|required_if=Type start"`
}

func NewChessGame(ws *websocket.Conn) *ChessGame {
	game := ChessGame{whiteWebsocket: ws}
	return &game
}

var ErrCannotJoinStartedGame = errors.New("cannot join a started game")

func (game *ChessGame) Join(ws *websocket.Conn) error {
	// you cannot join the same game twice
	if game.blackWebsocket != nil {
		return ErrCannotJoinStartedGame
	}
	game.blackWebsocket = ws
	whiteChannel := make(chan Message)
	blackChannel := make(chan Message)
	go playChess(game.whiteWebsocket, game.blackWebsocket, whiteChannel, blackChannel)
	go forwardFromWebsocketToChannel(game.whiteWebsocket, whiteChannel)
	go forwardFromWebsocketToChannel(game.blackWebsocket, blackChannel)
	return nil
}

const gameDurationInMinutes = 7
const incrementPerMoveInSeconds = 3

func playChess(
	whiteWebsocket, blackWebsocket *websocket.Conn,
	whiteChannel, blackChannel <-chan Message,
) {
	turnWhite := true
	whiteTimer := NewCountdown(gameDurationInMinutes*60, incrementPerMoveInSeconds, func() {})
	blackTimer := NewCountdown(gameDurationInMinutes*60, incrementPerMoveInSeconds, func() {})
	whiteWebsocket.WriteJSON(Message{Type: "start", Color: "white", WhiteTime: whiteTimer.GetRemaining(), BlackTime: blackTimer.GetRemaining()})
	blackWebsocket.WriteJSON(Message{Type: "start", Color: "black", WhiteTime: whiteTimer.GetRemaining(), BlackTime: blackTimer.GetRemaining()})
	whiteTimer.Start()
	for {
		select {
		case message := <-whiteChannel:
			switch message.Type {
			case "error":
				return
			case "move":
				if turnWhite {
					whiteTimer.Stop()
					blackTimer.Start()
					message.Color = "white"
					message.WhiteTime = whiteTimer.GetRemaining()
					message.BlackTime = blackTimer.GetRemaining()
					blackWebsocket.WriteJSON(message)
					whiteWebsocket.WriteJSON(message)
					turnWhite = false
				}
			}
		case message := <-blackChannel:
			switch message.Type {
			case "error":
				return
			case "move":
				if !turnWhite {
					blackTimer.Stop()
					whiteTimer.Start()
					message.Color = "black"
					message.WhiteTime = whiteTimer.GetRemaining()
					message.BlackTime = blackTimer.GetRemaining()
					whiteWebsocket.WriteJSON(message)
					blackWebsocket.WriteJSON(message)
					turnWhite = true
				}
			}
		}
	}
}

func forwardFromWebsocketToChannel(ws *websocket.Conn, ch chan<- Message) {
	defer ws.Close()
	for {
		message := Message{}
		err := ws.ReadJSON(&message)

		if err != nil {
			ch <- Message{Type: "error"}
			return
		}

		ch <- message
	}
}
