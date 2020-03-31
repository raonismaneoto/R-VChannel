package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

type RVMessage struct {
	Data string
}

type RVConfiguration struct {
	ChannelFilePath string
	PortToOpen string
	ServerIp string
	ClientIp string
	SudoPassword string
}

func getConfiguration(gopath string) RVConfiguration {
	log.Println("Starting reading configuration process")
	// it must open the port and make all scripts executable
	file, err := os.Open(gopath+"/src/github.com/raonismaneoto/R-VChannel/rvchannel.json")
	defer file.Close()

	if err != nil {
		log.Println("Error on opening configuration file", err.Error())
	}

	decoder := json.NewDecoder(file)
	configuration := RVConfiguration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Println("Error on decoding configuration file", err.Error())
	}

	return configuration
}

func setup(gopath string) RVConfiguration {
	// it must open the port and make all scripts executable
	configuration := getConfiguration(gopath)

	setupScriptPath := "$GOPATH/src/github.com/raonismaneoto/R-VChannel/setup.sh"

	exec.Command("/bin/sh", "-c", "chmod 777 " + setupScriptPath)

	channelScriptPath := "$GOPATH/src/github.com/raonismaneoto/R-VChannel/channel.sh"

	exec.Command("/bin/sh", "-c", "chmod 777 " + channelScriptPath)

	cmd := exec.Command("/bin/sh", "-c", setupScriptPath + " " + configuration.PortToOpen + " " +  configuration.ChannelFilePath)

	data, err := cmd.Output()

	if err != nil {
		log.Println("Error on executing setup script ", err.Error())
	}

	print(data)

	return configuration
}

func server(conChan chan string, configuration RVConfiguration) {
	log.Println("Starting server routine")
	l, err := net.Listen("tcp", configuration.ServerIp+":"+configuration.PortToOpen)

	if err != nil {
		log.Println("Error on opening listen connection", err.Error())
	}
	conChan <- "We can continue"
	c, err := l.Accept()

	if err != nil {
		log.Println("Error on accepting connection", err.Error())
	}

	for {
		log.Println("Starting server loop")
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

	gopath_cmd := exec.Command("/bin/sh", "-c", "echo $GOPATH")
	gopath, _ := gopath_cmd.Output()
	gopath_str := strings.TrimSpace(string(gopath))

	f, err := os.OpenFile(gopath_str+"/src/github.com/raonismaneoto/R-VChannel/logs/log.info", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	defer f.Close()
	wrt := io.MultiWriter(os.Stdout, f)
	log.SetOutput(wrt)

	configuration := setup(gopath_str)
	connChan := make(chan string)
	go server(connChan, configuration)
	<-connChan

	c := make(chan interface{})
	<-c
}
