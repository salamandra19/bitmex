package proto

// To get subscription client should send to the channel message
// JSON: {"action": "subscribe", "symbols": <[]string>}. Where
// field "symbols" is optional. If "symbols" is empty client will be
// subscribed on all symbols.
// To unsubscribe client should send to the channel message
// JSON: {"action": "unsubscribe"}.
// Data sending by websocket channel to the client will look like
// {
//    timestamp: <timestamp>,
//    symbol: <symbol_name>,
//    price: <lastPrice>
// }
type (
	MsgClientAction struct {
		Action  string   `json:"action"`
		Symbols []string `json:"symbols"`
	}

	MsgSrv struct {
		Timestamp string  `json:"timestamp"`
		Symbol    string  `json:"symbol"`
		Price     float64 `json:"lastPrice"`
	}
)
