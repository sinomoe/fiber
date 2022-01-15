package comet

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/sinomoe/fiber/pkg/dto/base"
	"time"
)

type Client struct {
	userId string
	comet  *Comet
	conn   *websocket.Conn
	Send   chan base.Message

	maxMessageSize int64
	pongWait       time.Duration
	pingPeriod     time.Duration
	writeWait      time.Duration
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
		c.comet.Unregister(c.userId)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			bs, _ := json.Marshal(message)
			w.Write(bs)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) heartbeat() {
	ticker := time.NewTicker(c.pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		c.comet.Unregister(c.userId)
	}()

	c.conn.SetReadLimit(c.maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(c.pongWait))
		return nil
	})

	for {
		select {
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(c.writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
