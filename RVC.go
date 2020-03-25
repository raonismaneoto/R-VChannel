package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"net"
	"os"
	"os/exec"
)

type RVMessage struct {
	Data string
}

//func setup() {
//	// it must open the port and make all scripts executable
//}

func server(conChan chan string) {
	l, _ := net.Listen("tcp", "192.168.15.15:3020")
	conChan <- "We can continue"
	for {
		print("receivedddd")
		c, err := l.Accept()


		message, err := bufio.NewReader(c).ReadBytes('\n')

		if err != nil {
			log.Fatalln("Error")
		}

		var msg RVMessage

		err = json.Unmarshal(message, &msg)
		fmt.Print(msg.Data)
		fmt.Print(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		cmdStr := "notify-send " + msg.Data
		cmd := exec.Command("/bin/sh", "-c", cmdStr)

		data, err := cmd.Output()

		if err != nil {
			println(err.Error())
			print(data)
			return
		}

	}
}

func client() {
	c, err := net.Dial("tcp", "192.168.15.15:3020")

	if err != nil {
		print("Error on connecting to server")
	}

	buffer := make(chan string, 10000000)


	go func() {
		for {
			watcher, _ := fsnotify.NewWatcher()
			// watch for error

			if err := watcher.Add("/home/raoni/tst"); err != nil {
				fmt.Println("ERROR", err)
			}

			select {
			// watch for events
			case <-watcher.Events:
				fmt.Printf("New event received")
				file, err := os.Open("/home/raoni/tst")
				if err != nil {
					fmt.Print("error on opening file")
				}

				scanner := bufio.NewScanner(file)
				print(scanner.Text())
				for scanner.Scan() {
					buffer <- scanner.Text()
					break
				}
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	for {
		message := &RVMessage{
			Data: <-buffer,
		}
		e, err := json.Marshal(message)

		if err != nil {
			log.Println("Error")
		}

		_, err = c.Write(append(e, '\n'))

		if err != nil {
			log.Println("Error when sending message")
		}
	}
	fmt.Print("dieing")
}

func main() {
	//setup()
	connChan := make(chan string)
	go server(connChan)
	<-connChan
	go client()

	c := make(chan interface{})
	<-c
}
