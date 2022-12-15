package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Session struct {
	clients map[*websocket.Conn]bool
	history []Message
	exited	chan int
	exit	chan int
}


func (s *Session) run(broadcast chan Message) {
	for {
		select {
		case msg := <-broadcast:
			for client := range s.clients {
				client.WriteJSON(msg)
			}
		case <-s.exit:
			for client := range s.clients {
				client.Close()
			}
			s.exited <- 1
			return
		}
	}
}

func makeSocket(w http.ResponseWriter, r *http.Request, s *Session, broadcast chan Message) {
		// upgrade http request to websocket
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		// add new websocket connection to clients
		s.clients[conn] = true
		defer func() {
			delete(s.clients, conn)
			conn.Close()
		}()

		// send history messages to the new connection
		for _, msg := range s.history {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("ERROR: %v", err)
				return
			}
		}

		// infinite loop for listening
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("ERROR: %v", err)
				return
			}

			broadcast <- msg
		}
}