package main

import (
	"fmt"
	"time"
)

type Message struct {
	From    string
	Payload string
}

type Server struct {
	msgch  chan Message
	quitch chan struct{}
}

func (s *Server) startAndListen() {
server:
	for {
		select {
		case msg := <-s.msgch:
			fmt.Printf("received message from: %s payload: %s\n", msg.From, msg.Payload)
		case <-s.quitch:
			fmt.Printf("Server shut down")
			break server
		default:
			fmt.Printf("empty channel\n")
			time.Sleep(2 * time.Second)
		}
	}
}

func sendMessageToServer(msgch chan Message, payload string) {
	msg := Message{
		From:    "Bob",
		Payload: payload,
	}
	msgch <- msg
	fmt.Printf("sending message from: %s\n", msg.From)
}

func shutDownServer(quitch chan struct{}) {
	close(quitch)
}

func main() {
	s := &Server{
		msgch:  make(chan Message),
		quitch: make(chan struct{}),
	}

	go s.startAndListen()
	for i := 0; i < 10; i++ {
		go func() {
			time.Sleep(1 * time.Second)
			sendMessageToServer(s.msgch, "Hello world!")
		}()
	}
	time.Sleep(15 * time.Second)
	shutDownServer(s.quitch)

	select {}
}
