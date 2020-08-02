package protocol

import (
	"fmt"
	"io"
)

type CommandWriter struct {
	writer io.Writer
}

func NewCommandWriter(writer io.Writer) *CommandWriter {
	return &CommandWriter{
		writer: writer,
	}
}

func (w *CommandWriter) writeString(msg string) error {
	_, err := w.writer.Write([]byte(msg))

	return err
}

func (w *CommandWriter) Write(command interface{}) error {
	// naive implementation ...
	var err error

	switch v := command.(type) {
	case SendCommand:
		err = w.writeString(fmt.Sprintf("SEND %v\n", v.Message))

	case MessageCommand:
		err = w.writeString(fmt.Sprintf("MESSAGE %v\t%v\t%v\t%v\n", v.Name, v.Message, v.Room, v.Date))

	case NameCommand:
		err = w.writeString(fmt.Sprintf("NAME %v\n", v.Name))

	case RoomCommand:
		err = w.writeString(fmt.Sprintf("ROOM %v\n", v.Room))

	default:
		err = UnknownCommand
	}

	return err
}
