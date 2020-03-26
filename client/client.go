package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"
)

type RVMessage struct {
	Data string
}

type RVConfiguration struct {
	ChannelFilePath string
	PortToOpen string
	ServerIp string
	ClientIp string
}

func getConfiguration() RVConfiguration {
	// it must open the port and make all scripts executable
	file, _ := os.Open("rvchannel.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := RVConfiguration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	return configuration
}

func client(configuration RVConfiguration) {
	c, err := net.Dial("tcp", configuration.ClientIp+":"+configuration.PortToOpen)
	if err != nil {
		flag := false
		for i := 0; i < 20; i++ {
			c, err = net.Dial("tcp", configuration.ClientIp+":"+configuration.PortToOpen)

			if err == nil {
				flag = true
				break
			}
			time.Sleep(2 * time.Minute)
		}
		if(!flag) {
			return
		}
	}

	buffer := make(chan string, 10000000)

	go func() {
		for {
			watcher, _ := fsnotify.NewWatcher()
			// watch for error

			if err := watcher.Add(configuration.ChannelFilePath); err != nil {
				fmt.Println("ERROR", err)
			}

			select {
			// watch for events
			case <-watcher.Events:
				fmt.Printf("New event received")
				body, err := ioutil.ReadFile(configuration.ChannelFilePath)

				if err != nil {
					fmt.Print("error on reading file")
				}

				buffer <- string(body)
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
	configuration := getConfiguration()
	go client(configuration)

	ThreadLocker := make(chan interface{})

	<- ThreadLocker
}
