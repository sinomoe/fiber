package comet

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/sinomoe/fiber/pkg/base"
	"log"
	"net/http"
	"sync"
	"time"
)

type Comet struct {
	address string
	clients map[string]*Client
	l       sync.RWMutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewComet(address string) *Comet {
	comet := &Comet{
		address: address,
		clients: make(map[string]*Client),
	}
	return comet
}

func (c *Comet) Spin() {
	http.HandleFunc("/ws", c.serveClient)
	http.ListenAndServe(c.address, nil)
}

func (c *Comet) serveClient(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user")
	if len(userId) == 0 {
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		userId:         userId,
		comet:          c,
		conn:           conn,
		Send:           make(chan base.Message, 16),
		maxMessageSize: 1024,
		pongWait:       60 * time.Second,
		pingPeriod:     50 * time.Second,
		writeWait:      10 * time.Second,
	}
	go client.heartbeat()
	go client.writePump()
	c.Register(userId, client)
}

func (c *Comet) Register(userId string, conn *Client) {
	c.l.Lock()
	defer c.l.Unlock()
	c.clients[userId] = conn
	return
}

func (c *Comet) Unregister(userId string) {
	c.l.Lock()
	defer c.l.Unlock()
	delete(c.clients, userId)
	return
}

func (c *Comet) GetClient(userId string) (*Client, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	if conn, ok := c.clients[userId]; ok {
		return conn, nil
	}
	return nil, errors.New("user not exist")
}
