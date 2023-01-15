package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type Session struct {
	clients 	map[*websocket.Conn]bool
	history 	[]Message
	broadcast 	chan Message
	jobs    	chan Message
	exited		chan int
	exit		chan int
}

const (
	nWorkers = 10
)

func (s *Session) run(gptClient *gogpt.Client, ctx context.Context) {
	exit := make(chan bool)
	
	// start session worker pool
	for i := 0; i < nWorkers; i++ {
		go startWorker(s.jobs, s.broadcast, gptClient, ctx, exit)
	}
	
	for {
		select {
		case msg := <-s.broadcast:
			s.history = append(s.history, msg)

			for client := range s.clients {
				client.WriteJSON(msg)
			}
		case <-s.exit:
			for client := range s.clients {
				client.Close()
			}
			for i := 0; i < nWorkers; i++ {
				exit <- true
			}
			s.exited <- 1
			return
		}
	}
}

func makeSocket(w http.ResponseWriter, 
				r *http.Request, 
				s *Session) {
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

			// push msg to jobs queue
			s.jobs <- msg
		}
}


func startWorker(jobs <-chan Message, broadcast chan<- Message, gptClient *gogpt.Client, ctx context.Context,
	exit <-chan bool) {
	gptChan := make(chan string)
	for {
		select {
		case job := <-jobs:
			if strings.HasPrefix(job.Text, "@GPT3") {
				req := gogpt.CompletionRequest{
					Model: "text-davinci-003",
					MaxTokens: 1000,
					Prompt: job.Text[5:],
				}
	
				go func() {
					res, err := gptClient.CreateCompletion(ctx, req)
					if err != nil {
						gptChan <- "Error:" + err.Error()
						return
					}
					gptChan <- res.Choices[0].Text
				}()
	
				broadcast <- job
				broadcast <- Message{Username: "GPT3", Text: <-gptChan}
			} else {
				broadcast <- job
			}

		case <-exit:
			return
		}
		
	}
}