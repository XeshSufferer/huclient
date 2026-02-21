package LittleMessenger

import (
	"fmt"

	"github.com/XeshSufferer/huclient"
	"github.com/gorilla/websocket"
)

func main() {
	// init
	client := huclient.NewClient("example.com", "app")

	client.On("ReceiveMessage", func(conn *websocket.Conn, message *huclient.Message) {
		fmt.Println("Receive Message!")
		fmt.Println(message.Args) // Write Raw JSON
	})

	client.Connect()

	for {
		var input string
		fmt.Scanln(&input)

		// Sending to server RPC "SendMessage" with arg = input
		err := client.SendMessage("SendMessage", input)

		if err != nil {
			fmt.Println(err)
		}
	}
}
