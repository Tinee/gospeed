package main

import (
	"log"

	"github.com/Tinee/gospeed/client"
)

type app struct {
}

func (a *app) HandleExit() {

}

func (a *app) HandleStart() {
	log.Println("START")
}

func (a *app) HandleSuggestion(b []byte) {

}

func main() {
	c, err := client.New("Marcus", ":8000")
	app := &app{}
	err = c.OpenAndListen(app)
	if err != nil {
		log.Fatalln(err)
	}

	c.Write("Hejsan")

	forever := make(chan bool)
	<-forever
}
