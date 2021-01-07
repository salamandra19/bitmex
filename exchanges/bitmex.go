package exchanges

import (
	"github.com/gorilla/websocket"
	"github.com/powerman/structlog"
	"github.com/salamandra19/bitmex/proto"
)

var (
	log = structlog.New()
)

type MsgBitmex struct {
	Table  string
	Action string
	Data   []MsgBitmexData
}

type MsgBitmexData struct {
	Timestamp string  `json:"timestamp"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"lastPrice"`
}

func NewBitmex(c chan proto.MsgSrv) {
	go func() {
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
				c <- convert(msg.Data[i])
			}
		}
	}()
}

func convert(data MsgBitmexData) proto.MsgSrv {
	return proto.MsgSrv{
		Timestamp: data.Timestamp,
		Symbol:    data.Symbol,
		Price:     data.Price,
	}
}
