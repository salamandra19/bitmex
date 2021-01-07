package exchanges

import (
	"github.com/gorilla/websocket"
	"github.com/powerman/structlog"
	"github.com/projects/bitmex/proto"
)

var (
	log = structlog.New()
)

type MsgBitmex struct {
	Table  string         `json:"table"`
	Action string         `json:"action"`
	Data   []proto.MsgSrv `json:"data"`
}

func NewBitmex(c chan proto.MsgSrv) {
	go func(c chan proto.MsgSrv) {
		conn, _, err := websocket.DefaultDialer.Dial("wss://testnet.bitmex.com/realtime?subscribe=instrument", nil)
		if err != nil {
			log.Fatal("failed to dial", "err", err)
		}
		defer conn.Close()

		for {
			var msg MsgBitmex
			err = conn.ReadJSON(&msg)
			if err != nil {
				log.Err("failed to read msg", "err", err)
				return
			}
			for i := range msg.Data {
				c <- msg.Data[i]
			}
		}
	}(c)
}
