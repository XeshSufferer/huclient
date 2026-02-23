package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/XeshSufferer/huclient"
	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("[client] starting...")

	client := huclient.NewClient("localhost:3000", "/app")

	client.OnConnected(func(c *huclient.Client) {
		fmt.Println("[client] OnConnected")
	})

	client.OnDisconnected(func(c *huclient.Client) {
		fmt.Println("[client] OnDisconnected")
	})

	if err := client.Connect(); err != nil {
		fmt.Printf("[client] connect error: %v\n", err)
		return
	}
	fmt.Println("[client] Connect() OK")

	client.On("pong", func(conn *websocket.Conn, msg *huclient.Message) {
		fmt.Printf("[client] PONG: %s\n", string(msg.Args))
	})

	fmt.Println("[client] press Enter to send ping")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println("[client] sending ping...")
		if err := client.SendMessage("ping", "hello"); err != nil {
			fmt.Printf("[client] send error: %v\n", err)
		}
	}
}
