package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nqbao/learn-go/chatserver/client"
	"github.com/nqbao/learn-go/chatserver/protocol"
)

var (
	address = flag.String("address", "localhost:3333", "address")
	debug   = flag.Bool("debug", false, "debug")
)

const (
	DELIMITER byte = '\n'
)

func main() {
	flag.Parse()
	//var c client.ChatClient
	c := client.NewClient()
	c.Dial(*address)
	defer c.Close()

	// start the client to listen for incoming message
	go c.Start()
	writer := c.GetCmdWriter()
	go func() {
		for {
			select {
			case err := <-c.Error():
				if err == io.EOF {
					fmt.Println("Connection closed connection from server.")
				} else {
					panic(err)
				}
			case msg := <-c.Incoming():
				// we need to make the change via ui update to make sure the ui is repaint correctly
				fmt.Printf("%v\t%v: %v\n", msg.Date, msg.Name, msg.Message)
			}
		}
	}()

	fmt.Printf("Enter your name: ")
	stdReader := bufio.NewReader(os.Stdin)
	name, _ := stdReader.ReadString('\n')
	c.SetName(name[:len(name)-1])
	fmt.Printf("Which chat room do you wanna join?\n")
	room, _ := stdReader.ReadString('\n')
	c.SetRoom(room[:len(room)-1])
	for {
		message, _ := stdReader.ReadString('\n')
		messageCmd := protocol.SendCommand{message[:len(message)-1]}
		if *debug {
			log.Printf("Sender: the request content: %v\n", messageCmd)
		}
		err := writer.Write(messageCmd)
		if err != nil {
			log.Printf("Sender: Write Error: %s\n", err)
		}
	}

}
