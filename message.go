package main

type Message struct {
	Type      string `json:"type"`
	Color     string `json:"color"`
	From      string `json:"from"`
	To        string `json:"to"`
	Promotion string `json:"promotion"`
}
