// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svc

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/powerman/structlog"
	"github.com/projects/bitmex/proto"
)

var (
	log = structlog.New()
)

type connection struct {
	ws      *websocket.Conn
	send    chan interface{}
	symbols []string
	mu      sync.Mutex
}

func newConnection(ws *websocket.Conn) *connection {
	return &connection{
		ws:   ws,
		send: make(chan interface{}, 256),
	}
}

func (c *connection) close() {
	close(c.send)
}

func (c *connection) reader() {
	for {
		var msg proto.MsgClientAction
		err := c.ws.ReadJSON(&msg)
		if err != nil {
			break
		}
		switch msg.Action {
		case "unsubscribe":
			return
		case "subscribe":
			if len(msg.Symbols) == 0 {
			} else {
				c.mu.Lock()
				for i := range msg.Symbols {
					c.symbols = append(c.symbols, msg.Symbols[i])
				}
				c.mu.Unlock()
			}
		default:
			log.Err("unsupported", "action", msg.Action)
			return
		}
	}

	c.ws.Close()
}

func (c *connection) writer() {
	for msg := range c.send {
		err := c.ws.WriteJSON(msg)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}