package websocket

import "github.com/gorilla/websocket"

var (
	newline = []byte{'\n'}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	server *Server
	conn   *websocket.Conn
	send   chan []byte
}

func (c *Client) read() {
	defer func() {
		c.server.unregister <- c
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.NextReader(); err != nil {
			c.conn.Close()
			break
		}
	}
}

func (c *Client) write() {
	for {
		select {
		case message := <-c.send:
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}
