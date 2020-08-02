package client

import (
	"log"
	"net"

	"github.com/nqbao/learn-go/chatserver/protocol"
)

// TcpChatClient ...
type TcpChatClient struct {
	conn      net.Conn
	cmdReader *protocol.CommandReader
	cmdWriter *protocol.CommandWriter
	name      string
	room      string
	error     chan error
	incoming  chan protocol.MessageCommand
}

// NewClient creates a new client with chanel or errors and message command
func NewClient() *TcpChatClient {
	return &TcpChatClient{
		incoming: make(chan protocol.MessageCommand),
		error:    make(chan error),
	}
}

// Dial create a connection with ther server at given address
func (c *TcpChatClient) Dial(address string) error {
	conn, err := net.Dial("tcp", address)

	if err == nil {
		c.conn = conn
		c.cmdReader = protocol.NewCommandReader(conn)
		c.cmdWriter = protocol.NewCommandWriter(conn)
	}

	return err
}
func (c *TcpChatClient) GetName() string {
	return c.name
}

func (c *TcpChatClient) GetConn() net.Conn {
	return c.conn
}

func (c *TcpChatClient) GetCmdReader() *protocol.CommandReader {
	return c.cmdReader
}

func (c *TcpChatClient) GetCmdWriter() *protocol.CommandWriter {
	return c.cmdWriter
}

// Start stay awake to read cmdReader
func (c *TcpChatClient) Start() {
	for {
		cmd, err := c.cmdReader.Read()

		if err != nil {
			c.error <- err
			break // TODO: find a way to recover from this
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.MessageCommand:
				c.incoming <- v
			default:
				log.Printf("Unknown command: %v", v)
			}
		}
	}
}

func (c *TcpChatClient) Close() {
	c.conn.Close()
}

func (c *TcpChatClient) Incoming() chan protocol.MessageCommand {
	return c.incoming
}

func (c *TcpChatClient) Error() chan error {
	return c.error
}

func (c *TcpChatClient) Send(command interface{}) error {
	return c.cmdWriter.Write(command)
}

func (c *TcpChatClient) SetName(name string) error {
	return c.Send(protocol.NameCommand{name})
}

func (c *TcpChatClient) SetRoom(name string) error {
	return c.Send(protocol.RoomCommand{name})
}

func (c *TcpChatClient) SendMessage(message string) error {
	return c.Send(protocol.SendCommand{
		Message: message,
	})
}
