package main

import (
	"log"
	"time"
)

type game struct {
	players   map[*player]bool
	sentences []sentence

	nextSentence chan *sentence
	register     chan *player
	unregister   chan *player
	answers      chan []byte
	shutdown     chan bool
	start        chan bool
}

func startGame() *game {
	g := &game{
		sentences: []sentence{
			{
				Content: "Testar detta lite grann",
			},
			{
				Content: "Testing",
			},
		},
		players:      make(map[*player]bool),
		register:     make(chan *player),
		unregister:   make(chan *player),
		nextSentence: make(chan *sentence),
		start:        make(chan bool),
		answers:      make(chan []byte),
		shutdown:     make(chan bool),
	}

	go g.track()

	return g
}

func (g *game) track() {
	for {
		select {
		case <-g.start:
			go g.sendFirstSentence()
		case p := <-g.register:
			g.players[p] = true
		case p := <-g.unregister:
			if _, ok := g.players[p]; ok {
				delete(g.players, p)
				// TODO
				// close(p.send)
			}
		case b := <-g.answers:
			proccessAnswer(b)
		case <-g.shutdown:
			log.Println("Closing the game")
			return
		}
	}
}

func (g *game) sendFirstSentence() {
	if len(g.sentences) == 0 {
		log.Fatalf("server doesn't have any sentences")
	}
	first := g.sentences[0]
	first.Sent = time.Now()

	g.nextSentence <- &first
}

func proccessAnswer(answer []byte) {
	log.Printf("Data is coming in %s", answer)
}
