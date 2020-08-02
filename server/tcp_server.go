package server

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nqbao/learn-go/chatserver/protocol"
)

type client struct {
	conn   net.Conn
	name   string
	room   string
	writer *protocol.CommandWriter
}

type TcpChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
	db       *sql.DB
}

var (
	UnknownClient = errors.New("Unknown client")
)

func NewServer() *TcpChatServer {
	return &TcpChatServer{
		mutex: &sync.Mutex{},
	}
}

func (s *TcpChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)
	if err == nil {
		s.listener = l
	}
	dns := fmt.Sprintf("thaovu:password@tcp(localhost:3306)/chatdb")
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	s.db = db
	log.Printf("Listening on %v", address)
	return err
}

func (s *TcpChatServer) Close() {
	s.listener.Close()
}

func (s *TcpChatServer) Start() {
	for {
		// XXX: need a way to break the loop
		conn, err := s.listener.Accept()

		if err != nil {
			log.Print(err)
		} else {
			// handle connection
			client := s.accept(conn)
			go s.serve(client)
		}
	}
}

func (s *TcpChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		// TODO: handle error here?
		switch v := command.(type) {
		case protocol.MessageCommand:
			if client.room == v.Room {
				client.writer.Write(command)
			}
		}
	}

	return nil
}

func (s *TcpChatServer) Send(name string, command interface{}) error {
	for _, client := range s.clients {
		if client.name == name {
			return client.writer.Write(command)
		}
	}

	return UnknownClient
}

func (s *TcpChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}

	s.clients = append(s.clients, client)

	return client
}

func (s *TcpChatServer) remove(client *client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// remove the connections from clients array
	for i, check := range s.clients {
		if check == client {
			s.clients = append(s.clients[:i], s.clients[i+1:]...)
		}
	}

	log.Printf("Closing connection from %v", client.conn.RemoteAddr().String())
	client.conn.Close()
}

func (s *TcpChatServer) insert(user string, text string, room string, date string) {
	stmtIns, err := s.db.Prepare("INSERT INTO chattable VALUES( ?, ?, ?, ? )") // ? = placeholder
	log.Printf("Insert values %s - %s - %s - %s \n", user, text, room, date)
	_, err = stmtIns.Query(user, text, room, date) // Insert tuples (i, i^2)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

func (s *TcpChatServer) loadHistory(conn net.Conn, room string) {
	cmdWriter := protocol.NewCommandWriter(conn)
	stmtIns, err := s.db.Prepare("SELECT * from chattable where room=? order by date") // ? = placeholder
	rows, err := stmtIns.Query(room)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		command := protocol.MessageCommand{
			Message: string(values[1]),
			Name:    string(values[0]),
			Room:    string(values[2]),
			Date:    string(values[3]),
		}
		cmdWriter.Write(command)

	}

}

func (s *TcpChatServer) serve(client *client) {

	cmdReader := protocol.NewCommandReader(client.conn)
	defer s.remove(client)

	for {
		cmd, err := cmdReader.Read()

		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
		}

		if cmd != nil {
			switch v := cmd.(type) {
			case protocol.SendCommand:
				log.Printf("Broadcasting message %s from user %s\n", v.Message, client.name)
				t := time.Now().Format("2006-01-02 15:04:05")
				s.insert(client.name, v.Message, client.room, t)
				go s.Broadcast(protocol.MessageCommand{
					Message: v.Message,
					Name:    client.name,
					Room:    client.room,
					Date:    t,
				})

			case protocol.NameCommand:
				log.Printf("Set name %s for connection %v\n", v.Name, client.conn.RemoteAddr().String())
				client.name = v.Name

			case protocol.RoomCommand:
				log.Printf("Set room %s for connection %v\n", v.Room, client.conn.RemoteAddr().String())
				client.room = v.Room
				s.loadHistory(client.conn, v.Room)
			}

		}

		if err == io.EOF {
			break
		}
	}
}
