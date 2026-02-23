# huclient

A Go WebSocket client for interacting with [husocket](https://github.com/XeshSufferer/husocket)-based servers.

## Description

`huclient` is a minimalist client for connecting to WebSocket servers. Features include:

- Server connection management
- Registration of incoming message handlers
- JSON-serialized message sending
- Concurrent-safe operation

## Installation

```bash
go get github.com/XeshSufferer/huclient
```

## Dependencies

- `github.com/gorilla/websocket` — WebSocket client library

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/XeshSufferer/huclient"
    "github.com/gorilla/websocket"
)

func main() {
    client := huclient.NewClient("localhost:3000", "/ws")

    // Register handler
    client.On("server_event", func(conn *websocket.Conn, msg *huclient.Message) {
        fmt.Printf("Received: %s\n", msg.Args)
    })

    // Connect
    if err := client.Connect(); err != nil {
        panic(err)
    }

    // Send a message
    err := client.SendMessage("client_event", map[string]string{
        "text": "Hello, server!",
    })
    if err != nil {
        fmt.Println("Send error:", err)
    }

    // Block to keep the program running
    select {}
}
```

## API

### Client

| Method | Description |
|--------|-------------|
| `NewClient(host, path string) *Client` | Create a new client instance |
| `Connect() error` | Establish WebSocket connection |
| `SendMessage(method string, args interface{}) error` | Send a JSON-serialized message |
| `On(method string, handler func(*websocket.Conn, *Message))` | Register a message handler |

### Message

```go
type Message struct {
    Method string          `json:"method"`
    Args   json.RawMessage `json:"args"`
}
```

## Message Format

Messages are sent and received in the following JSON format:

```json
{
  "method": "event_name",
  "args": { ... }
}
```

## Examples

See `_examples/LittleMessenger/` for a console chat application example.

## License

MIT
