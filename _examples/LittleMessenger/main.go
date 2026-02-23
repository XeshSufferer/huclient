package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/XeshSufferer/huclient"
	"github.com/gorilla/websocket"
)

func main() {
	// init
	client := huclient.NewClient("localhost:3000", "app")

	client.On("ReceiveMessage", func(conn *websocket.Conn, message *huclient.Message) {
		fmt.Println(string(message.Args)) // Write Raw JSON
	})

	client.On("OnJoin", func(conn *websocket.Conn, message *huclient.Message) {
		fmt.Println(string(message.Args))
	})

	client.On("OnLeave", func(conn *websocket.Conn, message *huclient.Message) {
		fmt.Println(string(message.Args))
	})

	client.Connect()
	fmt.Print("Enter nickname: ")
	var nickname string

	fmt.Scanln(&nickname)
	err := client.SendMessage("Join", nickname)

	if err != nil {
		fmt.Println(err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		// Sending to server RPC "SendMessage" with arg = input
		err := client.SendMessage("SendMessage", scanner.Text())

		if err != nil {
			fmt.Println(err)
		}
	}
}
