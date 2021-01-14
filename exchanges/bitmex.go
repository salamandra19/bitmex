package exchanges

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/powerman/must"
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

type bitmexOp struct {
	Op   string
	Args []string
}

const (
	verb      = "GET"
	apiKey    = "ORqVaoVf1TJrVnKexpWjHfjk"
	apiSecret = "mvK7p-zYF5He2eistXxXUvASoJWRGvp6eOO5TF2gn4BHI2iB"
	ws_URL    = "wss://testnet.bitmex.com/realtime"
)

// TODO must be reconnected with an exponential delay in case of connection error.

// NewBitmex gets websocket connection to Bitmex and receives change massages.
func NewBitmex(c chan proto.MsgSrv) {
	go func() {
		conn, _, err := websocket.DefaultDialer.Dial(subscribe(), makeHeader())
		must.NoErr(err)
		defer terminate(conn)

		for {
			var msg MsgBitmex
			err = conn.ReadJSON(&msg)
			if err != nil {
				log.PrintErr("failed to read msg", "err", err)
				return
			}
			for i := range msg.Data {
				c <- convert(msg.Data[i])
			}
		}
	}()
}

func subscribe() string {
	params := url.Values{
		"subscribe": []string{"instrument"},
	}
	return ws_URL + "?" + params.Encode()
}

func makeHeader() http.Header {
	u, err := url.Parse(ws_URL)
	must.NoErr(err)
	expires := strconv.Itoa(int(time.Now().Add(time.Minute).Round(time.Second).Unix()))
	data := []byte(verb + u.Path + expires)

	h := hmac.New(sha256.New, []byte(apiSecret))
	_, err = h.Write(data)
	must.NoErr(err)
	signature := hex.EncodeToString(h.Sum(nil))

	return http.Header{
		"api-expires":   []string{expires},
		"api-key":       []string{apiKey},
		"api-signature": []string{signature},
	}
}

func terminate(conn *websocket.Conn) {
	err := conn.WriteJSON(bitmexOp{
		Op:   "unsubscribe",
		Args: []string{"instrument"},
	})
	if err != nil {
		log.PrintErr("failed to unsubscribe", "err", err)
	}
	conn.Close()
}

func convert(data MsgBitmexData) proto.MsgSrv {
	return proto.MsgSrv{
		Timestamp: data.Timestamp,
		Symbol:    data.Symbol,
		Price:     data.Price,
	}
}
