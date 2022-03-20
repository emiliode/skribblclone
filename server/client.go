package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID      string          `json:"id"`
	Score   int             `json:"score"`
	Conn    *websocket.Conn `json:"-"`
	Game    *Game           `json:"-"`
	Creator bool            `json:"creator"`
}
type EventMessage struct {
	Event  string  `json:"event"`
	Client *Client `json:"client"`
	Body   string  `json:"body"`
}
type Message struct {
	Event string `json:"event"`
	Body  string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Game.Leave <- c
		c.Conn.Close()
	}()
	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		var received_data Message
		json.Unmarshal(p, &received_data)
		fmt.Printf("JSON Received %+v", received_data)
		switch received_data.Event {
		case "GUESS":
			fmt.Printf("\033[32m Received GUESS %v from: %v \033[0m \n", received_data.Body, c.ID)
			c.Game.Broadcast <- EventMessage{Event: "GUESS", Body: received_data.Body, Client: c}

		default:
			message := Drawing{Client: c, Content: string(p)}
			c.Game.Draw <- message
			fmt.Printf("Message Received:%d %+v\n", messageType, message)

		}
	}
}
