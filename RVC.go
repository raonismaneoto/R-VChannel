package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func setup() RVConfiguration {
	// it must open the port and make all scripts executable
	configuration := getConfiguration()

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

func main() {
	configuration := setup()
	connChan := make(chan string)
	go server(connChan, configuration)
	<-connChan

	c := make(chan interface{})
	<-c
}
