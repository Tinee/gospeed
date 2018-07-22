package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

type player struct {
	conn      *websocket.Conn
	game      *game
	sentences []sentence
	name      string
}

func (p *player) addSentence(s sentence) {
	s.Arrived = time.Now()
	p.sentences = append(p.sentences, s)
}

func (p *player) readPump() {
	defer func() {
		p.game.unregister <- p
		p.conn.Close()
	}()

	p.conn.SetReadLimit(maxMessageSize)
	p.conn.SetReadDeadline(time.Now().Add(pongWait))
	p.conn.SetPongHandler(func(string) error {
		p.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := p.conn.ReadMessage()
		fmt.Println(message)
		if err != nil {
			break
		}

		p.game.start <- true
	}
}

func (p *player) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.conn.Close()
	}()

	for {
		select {
		case sentence, ok := <-p.game.nextSentence:
			if !ok {
				p.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			fmt.Println(sentence)
		case <-p.game.start:
			err := p.conn.WriteMessage(50, []byte(""))
			if err != nil {
				log.Print(err)
			}
			return
		case <-ticker.C:
			p.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
