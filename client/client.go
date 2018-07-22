package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// Client represents the client.
type Client struct {
	URL  string
	name string
	conn *websocket.Conn

	GoSpeedHandler
}

type GoSpeedHandler interface {
	HandleSuggestion([]byte)
	HandleExit()
	HandleStart()
}

var (
	// ErrInvalidArgument gets returned if the user passes in a wrong argument
	ErrInvalidArgument      = errors.New("invalid arguments")
	defaultHandleSuggestion = func([]byte) {}
)

// New creates a new Client.
func New(name, port string) (*Client, error) {
	if name == "" {
		return nil, ErrInvalidArgument
	}

	URL := url.URL{
		Scheme: "ws",
		Host:   port,
		Path:   "/",
	}

	client := &Client{
		URL:  URL.String(),
		name: name,
	}

	return client, nil
}

func (c *Client) Write(s string) {
	c.conn.WriteMessage(1, []byte(s))
}

// OpenAndListen attempts to connect to an websockets server.
func (c *Client) OpenAndListen(goh GoSpeedHandler) error {
	header := c.getAuthHeaders()

	conn, _, err := websocket.DefaultDialer.Dial(c.URL, header)
	if err != nil {
		return err
	}
	c.conn = conn
	c.GoSpeedHandler = goh

	go c.listen()
	return nil
}

// Close closes the clients connection
func (c *Client) Close() {
	c.Close()
}

func (c *Client) listen() error {
	for {
		messageType, b, err := c.conn.ReadMessage()
		if err != nil {
			return err
		}
		fmt.Print(messageType)
		switch messageType {
		case websocket.TextMessage:
			c.HandleSuggestion(b)
		case 50:
			c.HandleStart()
		}
	}
}

func (c *Client) getAuthHeaders() http.Header {
	header := make(http.Header)
	header["Authorization"] = []string{c.name}

	return header
}
