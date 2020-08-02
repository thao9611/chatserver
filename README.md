# Simple Chat Server

A simple chat server written in Golang, with very basic features:

  * Multiple chat rooms with historical data saved in mysql
  * User can connect to the chat server
  * User can set their name and choose to enter the chat room they want
  * User can send message to the chat room

## Protocol

For this excersie , I will use simple text-based message over TCP:

  * All messages are terminated with `\n`
  * To send a chat message, client will send: 
    * `SEND chat message`
    * For now, chat message can not contain new line.
  * To set client name, client will send:
    * `NAME username`
    * For now, username can not contain space
  * To enter a chat room, client will send:
    * `ROOM room`
    * The history of chat room will be loaded after room setting. Users are allowed to switch between rooms during one login. 
  * Server will send the following command to all clients when there are new message:
    * `MESSAGE username the actual message room date`
    * The message will be also inserted to database



## References

  * https://gist.github.com/drewolson/3950226
  * https://scotch.io/bar-talk/build-a-realtime-chat-server-with-go-and-websockets
  * [tui-go](https://github.com/marcusolsson/tui-go) for chat client. tview may be a better option, but textare input is not available yet.
  * https://github.com/nqbao/go-sandbox/chatserver 
