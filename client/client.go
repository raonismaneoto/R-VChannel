package main

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
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

func getConfiguration(gopath string) RVConfiguration {
	log.Println("Getting Configuration")
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
		log.Println("Error on decoding file to variable", err.Error())
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

	configuration := getConfiguration(gopath_str)
	go client(configuration)

	ThreadLocker := make(chan interface{})

	<- ThreadLocker
}
