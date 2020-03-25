package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
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

func setup() RVConfiguration {
	// it must open the port and make all scripts executable
	file, _ := os.Open("rvchannel.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := RVConfiguration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
	}

	setupScriptPath := "$HOME/go/src/github.com/raonismaneoto/R-VChannel/setup.sh"

	exec.Command("/bin/sh", "-c", "chmod 777 " + setupScriptPath)

	channelScriptPath := "$HOME/go/src/github.com/raonismaneoto/R-VChannel/channel.sh"

	exec.Command("/bin/sh", "-c", "chmod 777 " + channelScriptPath)

	cmd := exec.Command("/bin/sh", "-c", setupScriptPath + " " + configuration.PortToOpen + " " +  configuration.ChannelFilePath)

	data, err := cmd.Output()

	if err != nil {
		print("Error on setup script", err.Error())
	}

	print(data)

	return configuration
}

func server(conChan chan string, configuration RVConfiguration) {
	l, _ := net.Listen("tcp", configuration.ServerIp+":"+configuration.PortToOpen)
	conChan <- "We can continue"
	c, _ := l.Accept()
	for {
		message, err := bufio.NewReader(c).ReadBytes('\n')

		if err != nil {
			log.Fatalln("Error")
		}

		var msg RVMessage

		err = json.Unmarshal(message, &msg)
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
	configuration := setup()
	connChan := make(chan string)
	go server(connChan, configuration)
	<-connChan
	go client(configuration)

	c := make(chan interface{})
	<-c
}
