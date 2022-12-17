package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	gogpt "github.com/sashabaranov/go-gpt3"
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
			s.history = append(s.history, msg)

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

func makeSocket(w http.ResponseWriter, 
				r *http.Request, 
				s *Session, 
				broadcast chan Message,
				gptClient *gogpt.Client, 
				ctx context.Context) {
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

		gptChan := make(chan string)

		// infinite loop for listening
		for {
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("ERROR: %v", err)
				return
			}

			fmt.Println(msg);

			if strings.HasPrefix(msg.Text, "@GPT3") {
				req := gogpt.CompletionRequest{
					Model: "text-davinci-003",
					MaxTokens: 1000,
					Prompt: msg.Text[5:],
				}
	
				go func() {
					res, err := gptClient.CreateCompletion(ctx, req)
					if err != nil {
						gptChan <- "Error:" + err.Error()
						return
					}
					gptChan <- res.Choices[0].Text
				}()

				broadcast <- msg
				broadcast <- Message{Username: "GPT3", Text: <-gptChan}
			} else {
				broadcast <- msg
			}
		}
}