package main

import (
	"flag"
	"fmt"
	"github.com/anthdm/hollywood/actor"
	"github.com/anthdm/hollywood/remote"
	"github.com/gorilla/websocket"
	"net/http"
)

type PlayerSession struct {
	clientID int
	username string
	inLobby  bool
	conn     *websocket.Conn
}

func (p PlayerSession) Receive(context *actor.Context) {
	//TODO implement me
	panic("implement me")
}

func NewPlayerSession(clientID int, username string, conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			clientID: clientID,
			username: username,
			conn:     conn,
			inLobby:  false,
		}
	}
}

type GameServer struct {
}

func newGameServer() actor.Receiver {
	return &GameServer{}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		s.startHTTP()
		_ = msg
	}
}

func (s *GameServer) startHTTP() {
	fmt.Println("starting HTTP server on port -> 40000")

	go func() {
		http.HandleFunc("/ws", s.handleWS)
		err := http.ListenAndServe("0.0.0.0:40000", nil)

		if err != nil {
			return
		}
	}()
}

// handles the upgrade of the websocket
func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new client trying to connect")
	upgrade := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrade.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("ws upgrade err: ", err)
		return
	}

	fmt.Print("new client trying to connect")
	fmt.Print(conn)
}

func main() {
	var (
		listenAt = flag.String("listen", "127.0.0.1:40000", "")
	)

	rem := remote.New(*listenAt, remote.NewConfig())
	e, err := actor.NewEngine(actor.NewEngineConfig().WithRemote(rem))

	if err != nil {
		panic(err)
	}

	e.Spawn(newGameServer, "server")
	select {}
}
