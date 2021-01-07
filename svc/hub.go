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
	"fmt"

	"github.com/projects/bitmex/proto"
)

type hub struct {
	connections map[*connection]bool
	broadcast   chan proto.MsgSrv
	register    chan *connection
	unregister  chan *connection
}

func newHub() *hub {
	h := &hub{
		connections: make(map[*connection]bool),
		broadcast:   make(chan proto.MsgSrv),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
	}
	go h.run()
	return h
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
			fmt.Println("registered", c, h.connections)
		case c := <-h.unregister:
			if h.connections[c] {
				delete(h.connections, c)
				c.close()
			}
			fmt.Println("unregistered", c, h.connections)
		case m := <-h.broadcast:
			for c := range h.connections {
				if filter(m, *c) {
					select {
					case c.send <- m:
					default:
						delete(h.connections, c)
						c.close()
					}
				}
			}
		}
	}
}

func filter(m proto.MsgSrv, c connection) bool {
	if len(c.symbols) == 0 {
		return true
	}
	for i := range c.symbols {
		if m.Symbol == c.symbols[i] {
			return true
		}
	}
	return false
}
