package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var games = make(map[string]*Game)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return conn, nil
}

func serveWs(game *Game, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket endpoint hit")
	conn, err := Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
	}
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
    fmt.Println("Received Name %s",string(p))
    creator:=false
	for client, _ := range game.Clients {
		if client.ID == string(p) {
			if client.Conn != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("NAMEALREADYINUSE"))
				conn.Close()
                fmt.Printf("NAMEALREADYINUSE %+v", client)
				return
			}else{
                delete( game.Clients, client)
                creator= true
                continue
            }
		}
	}
    client := &Client{Conn: conn, Game: game, ID: string(p), Score: 0,Creator:creator}
    fmt.Printf("CLient Created %+v", client)
	game.Join <- client
	client.Read()

}
func handleCreateGame(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	client := &Client{Conn: nil, Game: nil, ID: string(b), Score: 0}
	game := NewGame(client)
	games[game.ID] = game
    w.Header().Set("Content-Type","application/json")
	fmt.Fprintf(w, "{ \"game\": \""+game.ID+"\" }")
    go game.Run() 

}
func setupRoutes() {
	http.HandleFunc("/game/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("GAME HIT")
		serveWs(games[r.URL.Path[len("/game/"):]], w, r)
	})
	http.HandleFunc("/creategame", handleCreateGame)
}

func main() {
	fmt.Println("Starting server for skribblclone")
	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
