package main

import (
	"github.com/gin-gonic/gin"
	"github.com/salamandra19/bitmex/exchanges"
	"github.com/salamandra19/bitmex/proto"
	"github.com/salamandra19/bitmex/svc"
)

func main() {
	var msgSrvChan = make(chan proto.MsgSrv)

	exchanges.NewBitmex(msgSrvChan)

	server := svc.NewServer(msgSrvChan)

	r := gin.Default()
	r.GET("/bitmex", func(c *gin.Context) {
		server.Handler(c.Writer, c.Request)
	})
	r.Run(":8844")
}
