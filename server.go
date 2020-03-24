package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type RVMessage struct{
	Data string
}

func main() {
	l, err := net.Listen("tcp", "192.168.15.11:3020")
	c, err := l.Accept()

	message, err := bufio.NewReader(c).ReadBytes('\n')

	if err != nil {
		log.Fatalln("Error")
	}

	var msg RVMessage

	err = json.Unmarshal(message, &msg)

	fmt.Printf("Msg %s received\n", msg.Data)
}
