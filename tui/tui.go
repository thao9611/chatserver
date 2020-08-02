package tui

import (
	"fmt"
	"io"

	"github.com/marcusolsson/tui-go"
	"github.com/nqbao/learn-go/chatserver/client"
)

func StartUi(c client.ChatClient) {
	loginView := NewLoginView()
	chatView := NewChatView()
	roomView := NewRoomView()

	ui, err := tui.New(loginView)
	if err != nil {
		panic(err)
	}

	quit := func() { ui.Quit() }
	changeRoom := func() { ui.SetWidget(roomView) }

	ui.SetKeybinding("Esc", changeRoom)
	ui.SetKeybinding("Ctrl+c", quit)

	//set handler for login
	loginView.OnLogin(func(username string) {
		c.SetName(username)
		ui.SetWidget(roomView)
	})
	roomView.OnRoom(func(room string) {
		c.SetRoom(room)
		ui.SetWidget(chatView)
	})

	//set handler for sending message
	chatView.OnSubmit(func(msg string) {
		c.SendMessage(msg)
	})

	go func() {
		for {
			select {
			case err := <-c.Error():

				if err == io.EOF {
					ui.Update(func() {
						chatView.AddMessage("Connection closed connection from server.")
					})
				} else {
					panic(err)
				}
			case msg := <-c.Incoming():
				// we need to make the change via ui update to make sure the ui is repaint correctly
				ui.Update(func() {
					chatView.AddMessage(fmt.Sprintf("%v  %v: %v\n", msg.Date, msg.Name, msg.Message))
				})
			}
		}
	}()

	if err := ui.Run(); err != nil {
		panic(err)
	}
}
