package main

import (
	"fmt"
	"github.com/edfcsx/gameserver/types"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"math/rand"
	"time"
)

const wsServerEndpoint = "ws://localhost:40000/ws"

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	username string
}

func (c *GameClient) login() error {
	return c.conn.WriteJSON(types.Login{
		ClientID: c.clientID,
		Username: c.username,
	})
}

func NewGameClient(conn *websocket.Conn, username string) *GameClient {
	return &GameClient{
		conn:     conn,
		clientID: rand.Intn(math.MaxInt),
		username: username,
	}
}

func main() {
	fmt.Println("Hello, client!")
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	fmt.Println("Hello, client!1")

	dialer.HandshakeTimeout = time.Second * 2
	conn, _, err := dialer.Dial(wsServerEndpoint, nil)

	fmt.Println("Hello, client!2")

	if err != nil {
		log.Fatalln(err)
	}

	print("Connected to server")

	c := NewGameClient(conn, "edfcsx")
	if err := c.login(); err != nil {
		log.Fatal(err)
	}
}
