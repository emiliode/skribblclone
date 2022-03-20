package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"skribblclone/server/internal/utils"
)

// TODO: load WORDLIST from file
var WORDLIST = [...]string{"Zimmer", "Wohnzimmer", "Schlafzimmer", "Kinderzimmer", "Badezimmer", "Küche", "Treppenhaus", "Treppe", "Teppich", "Vorhang", "Klingel", "klopfen", "kochen", "Topf", "Teller", "Besteck", "Gabel", "Messer", "Löffel", "Schere", "Lichtschalter", "Lampe", "Wand", "Fußboden", "Dach", "Kamin", "Heizung", "Spielzeug", "Schreibtisch", "Schrank", "Regal", "Bett", "Sofa", "Fernseher", "Radio", "Fernbedienung", "kochen", "Spühlmaschine", "Waschmaschine", "Trockner", "Föhn", "Rasierer", "Dusche", "Toilette", "Stuhl", "Tisch"}

type Game struct {
	Join        chan *Client
	Leave       chan *Client
	Draw        chan Drawing
	Clients     map[*Client]bool
	CurrentWord string
	ID          string
	Broadcast   chan EventMessage
}
type Drawing struct {
	Client  *Client `json:"client"`
	Content string  `json:"msg"`
}

func NewGame(client *Client) *Game {
	game := Game{
		Join:      make(chan *Client),
		Leave:     make(chan *Client),
		Clients:   make(map[*Client]bool),
		Broadcast: make(chan EventMessage),
		Draw:      make(chan Drawing),
		ID:        utils.GenerateGameID(10),
	}
	game.Clients[client] = true
    client.Creator = true
	return &game
}
func (game *Game) Run() {
	for {
		select {
		case client := <-game.Join:
			game.Clients[client] = true
			clients := []*Client{}
			for key, _ := range game.Clients {
				clients = append(clients, key)
			}
			data, err := json.Marshal(clients)
			if err != nil {
				fmt.Println(err)
				return
			}
			client.Conn.WriteJSON(Message{Event: "JOINCOMPLETE", Body: string(data)})
			fmt.Println("Size of Connection Pool: ", len(game.Clients))
			sendall_except_sender(&EventMessage{Event: "join", Client: client}, &game.Clients, client)
			break
		case client := <-game.Leave:
			delete(game.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(game.Clients))

			sendall(&EventMessage{Event: "leave", Client: client}, &game.Clients)

			break
		case drawing := <-game.Draw:
			fmt.Println("Drawing from " + drawing.Client.ID)
			for client, _ := range game.Clients {

				if client.ID == drawing.Client.ID {
					continue
				}
				fmt.Println(client)
				client.Conn.WriteMessage(websocket.TextMessage, []byte("{ \"event\": \"DRAWUPDATE\",\"body\": "+drawing.Content+" }"))
			}
			break
		case message := <-game.Broadcast:
			fmt.Println("Sending Message to all clients in Pool")
			sendall_except_sender(&message, &game.Clients, message.Client)
			break
		}
	}
}
func sendall(message *EventMessage, clients *map[*Client]bool) {
	for client, _ := range *clients {
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
			return

		}

	}
}
func sendall_except_sender(message *EventMessage, clients *map[*Client]bool, sender *Client) {

	for client, _ := range *clients {
		if client.ID == sender.ID {
			continue
		}
		if err := client.Conn.WriteJSON(message); err != nil {
			fmt.Println(err)
			return
		}
	}
}
func setnextword(g *Game) {
	pick := rand.Intn(len(WORDLIST))
	g.CurrentWord = WORDLIST[pick]
}
