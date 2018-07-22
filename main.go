package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type websocketServer struct {
	game     *game
	upgrader websocket.Upgrader
}

func newWebsocketServer(g *game) *websocketServer {
	return &websocketServer{
		game: g,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (ws *websocketServer) serveWebSocket(w http.ResponseWriter, r *http.Request) {
	playerName := r.Header["Authorization"][0]
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	player := &player{
		conn: conn,
		game: ws.game,
		name: playerName,
	}

	go func() {
		ws.game.register <- player
	}()
	go player.readPump()
	go player.writePump()
}

func main() {
	game := startGame()
	server := newWebsocketServer(game)

	http.HandleFunc("/", server.serveWebSocket)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
