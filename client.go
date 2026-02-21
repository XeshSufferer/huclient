package huclient

import (
	"encoding/json"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Method string          `json:"method"`
	Args   json.RawMessage `json:"args"`
}

type Client struct {
	Host     string
	Path     string
	Conn     *websocket.Conn
	handlers map[string]func(*websocket.Conn, *Message)
	mu       sync.Mutex
}

func NewClient(host, path string) *Client {
	return &Client{Host: host, Path: path, handlers: make(map[string]func(*websocket.Conn, *Message)), mu: sync.Mutex{}}
}

func (c *Client) Connect() error {
	connUrl := url.URL{Scheme: "ws", Host: c.Host, Path: c.Path}
	ws, _, err := websocket.DefaultDialer.Dial(connUrl.String(), nil)

	if err != nil {
		return err
	}
	c.Conn = ws

	go func() {
		for {
			mt, message, readErr := c.Conn.ReadMessage()
			if readErr != nil {
				continue
			}

			if mt != websocket.TextMessage {
				continue
			}

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			if h, ok := c.handlers[msg.Method]; ok {
				h(ws, &msg)
			}
		}
	}()
	return nil
}

func (c *Client) SendMessage(method string, args interface{}) error {

	rawArgs, err := json.Marshal(args)
	if err != nil {
		return err
	}

	content, err := json.Marshal(Message{Method: method, Args: rawArgs})

	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Conn.WriteMessage(websocket.TextMessage, content)
	return nil
}

func (c *Client) On(method string, f func(*websocket.Conn, *Message)) {
	c.mu.Lock()
	c.handlers[method] = f
	c.mu.Unlock()
}
