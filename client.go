package huclient

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Method string          `json:"method"`
	Args   json.RawMessage `json:"args"`
}

type Client struct {
	Host           string
	Path           string
	Conn           *websocket.Conn
	handlers       map[string]func(*websocket.Conn, *Message)
	mu             sync.Mutex
	onConnected    func(conn *Client)
	onDisconnected func(c *Client)
}

func NewClient(host, path string) *Client {
	return &Client{Host: host, Path: path, handlers: make(map[string]func(*websocket.Conn, *Message)), mu: sync.Mutex{}}
}

func (c *Client) OnConnected(f func(c *Client)) {
	c.onConnected = f
}

func (c *Client) OnDisconnected(f func(c *Client)) {
	c.onDisconnected = f
}

func (c *Client) Connect() error {
	log.SetPrefix("[huclient] ")
	connUrl := url.URL{Scheme: "ws", Host: c.Host, Path: c.Path}
	ws, _, err := websocket.DefaultDialer.Dial(connUrl.String(), nil)

	if err != nil {
		return err
	}
	c.Conn = ws

	c.On("close", func(conn *websocket.Conn, message *Message) {
		log.Printf("CONNECTION CLOSED: %s", string(message.Args))
	})

	go func() {
		closed := false
		for {
			if closed {
				return
			}

			mt, message, readErr := c.Conn.ReadMessage()

			log.Printf("RAW: len = %d | data = %s | mt = %d ", len(message), string(message), mt)
			if readErr != nil {
				if closeErr, ok := readErr.(*websocket.CloseError); ok {
					log.Printf("PROTOCOL CLOSE: code=%d, text=%s", closeErr.Code, closeErr.Text)
					if c.onConnected != nil {
						c.onDisconnected(c)
					}
				} else {
					log.Printf("READ ERR: %v", readErr)
				}
				return
			}

			if mt != websocket.TextMessage {
				continue
			}

			var msg Message
			if err := json.Unmarshal(message, &msg); err != nil {
				continue
			}

			log.Println(string(message))
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
	log.Printf("SENDING: %s", string(content))
	c.Conn.WriteMessage(websocket.TextMessage, content)
	return nil
}

func (c *Client) On(method string, f func(*websocket.Conn, *Message)) {
	c.mu.Lock()
	c.handlers[method] = f
	c.mu.Unlock()
}

func (c *Client) Close() {
	c.SendMessage("close", "")
	c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}
