package main

type Message struct {
	Type      string `json:"type" validate:"required,oneof=start move"`
	Color     string `json:"color" validate:"oneof=white black,required_if=Type start"`
	From      string `json:"from" validate:"required_if=Type move"`
	To        string `json:"to" validate:"required_if=Type move"`
	Promotion string `json:"promotion" validate:"oneof=q r b k,required_if=Type move"`
}
