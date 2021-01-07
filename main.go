package main

import (
	"log"
	"net/http"

	"github.com/salamandra19/bitmex/exchanges"
	"github.com/salamandra19/bitmex/proto"
	"github.com/salamandra19/bitmex/svc"
)

func main() {
	var msgSrvChan = make(chan proto.MsgSrv)

	exchanges.NewBitmex(msgSrvChan)

	server := svc.NewServer(msgSrvChan)
	http.Handle("/bitmex", http.HandlerFunc(server.Handler))
	log.Fatal(http.ListenAndServe(":8844", nil))
}
