package svc

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/salamandra19/bitmex/proto"
)

type Server struct {
	upgrader   *websocket.Upgrader
	hub        *hub
	msgSrvChan chan proto.MsgSrv
}

func NewServer(msgSrvChan chan proto.MsgSrv) *Server {
	srv := &Server{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
		hub:        newHub(),
		msgSrvChan: msgSrvChan,
	}
	go srv.sendingMsgSrv()
	return srv
}

// Handler is a HandlerFunc handling the client websocket interaction.
func (srv *Server) Handler(w http.ResponseWriter, r *http.Request) {
	ws, err := srv.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := newConnection(ws)

	srv.hub.register <- c
	defer func() { srv.hub.unregister <- c }()
	go c.writer()
	c.reader()
}

func (srv *Server) sendingMsgSrv() {
	for msg := range srv.msgSrvChan {
		srv.hub.broadcast <- msg
	}
}
