package main

import (
	"time"

	"github.com/XeshSufferer/huclient"
)

func main() {
	client := huclient.NewClient("localhost:3000", "app")
	client.Connect()
	//client.SendMessage("close_now", "")
	//select {}
	time.Sleep(1 * time.Second)
	client.Close()
}
