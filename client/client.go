package main

import (
	"encoding/json"
	"log"
	"net"
)

type RVMessage struct{
	Data string
}

func main() {
	c, err := net.Dial("tcp", "192.168.15.11:3020")

	message := &RVMessage{
		Data: "Oiiiii",
	}

	e, err := json.Marshal(message)

	if err != nil {
		log.Println("Error")
	}

	_, err = c.Write(append(e, '\n'))

	if err != nil {
		log.Println("Error")
	}
}
