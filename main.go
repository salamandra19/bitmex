package main

import (
	"flag"

	"github.com/gin-gonic/gin"
	"github.com/powerman/must"
	"github.com/salamandra19/bitmex/exchanges"
	"github.com/salamandra19/bitmex/proto"
	"github.com/salamandra19/bitmex/svc"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8844", "serve address")
	flag.Parse()

	var msgSrvChan = make(chan proto.MsgSrv)
	exchanges.NewBitmex(msgSrvChan)

	server := svc.NewServer(msgSrvChan)

	r := gin.Default()
	r.GET("/bitmex", func(c *gin.Context) {
		server.Handler(c.Writer, c.Request)
	})
	must.NoErr(r.Run(*addr))
}
