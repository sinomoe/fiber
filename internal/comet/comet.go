package comet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sinomoe/fiber/internal/logic/dto"

	"github.com/sinomoe/fiber/internal/config"

	"github.com/gorilla/websocket"
	"github.com/sinomoe/fiber/pkg/dto/base"
)

type Comet struct {
	address string
	clients map[string]*Client
	l       sync.RWMutex

	logicCli *http.Client
	cfg      *config.Comet
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewComet(cfg *config.Comet) *Comet {
	comet := &Comet{
		address:  fmt.Sprintf(":%d", cfg.WebsocketPort),
		clients:  make(map[string]*Client),
		logicCli: &http.Client{},
		cfg:      cfg,
	}
	return comet
}

func (c *Comet) Spin() {
	http.HandleFunc("/ws", c.serveClient)
	http.ListenAndServe(c.address, nil)
}

func (c *Comet) serveClient(w http.ResponseWriter, r *http.Request) {
	var (
		token = r.URL.Query().Get("token")
		user  string
		conn  *websocket.Conn
		err   error
	)
	if user, err = c.validateToken(token); err != nil {
		log.Println(err)
		return
	}
	if conn, err = upgrader.Upgrade(w, r, nil); err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		userId:         user,
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
	c.Register(user, client)
}

func (c *Comet) validateToken(token string) (username string, err error) {
	var (
		req  *http.Request
		resp *http.Response
		bs   []byte
	)
	if bs, err = json.Marshal(dto.ParseTokenRequest{Token: token}); err != nil {
		return
	}
	if req, err = http.NewRequest(http.MethodPost, c.cfg.LogicUrl, bytes.NewReader(bs)); err != nil {
		return
	}
	if resp, err = c.logicCli.Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = errors.New("token invalid")
		return
	}
	var r dto.ParseTokenResponse
	if err = json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return
	}
	username = r.Username
	return
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
