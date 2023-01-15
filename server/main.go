package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type Message struct {
	Username string `json:"username"`
	Text string `json:"text"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	clients = make(map[*websocket.Conn]bool)
)

const (
	apiKey = "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"
)


func runServer(sig chan os.Signal, exit chan int) {
	session := Session{
		clients: make(map[*websocket.Conn]bool),
		broadcast: make(chan Message),
		jobs: make(chan Message),
		exited: make(chan int),
		exit: make(chan int),
	}
	router := mux.NewRouter()

	// create gpt3 client
	ctx := context.Background()
	gptClient := gogpt.NewClient(apiKey)

	router.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Got connection")
		makeSocket(w, r, &session)
	})

	handler := cors.AllowAll().Handler(router)

	srv := &http.Server{
		Handler: handler,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go func() { log.Fatal(srv.ListenAndServe()) }()

	// start goroutine for broadcasting messages
	go session.run(gptClient, ctx)

	fmt.Println("Server listening on", srv.Addr)

	minutes := 0
	for { // handles gracefully exiting on signal
		select {
		case <-sig: // if there is an exit signal
			session.exit <- 1
			<-session.exited
			exit <- 1
			return
		case <-time.After(time.Minute * 1): // log how long the server has been running
			minutes += 1
			fmt.Printf("%d minutes have passed\n", minutes*1)
		}
	}

}


func main() {
	sigChan := make(chan os.Signal)
	exitChan := make(chan int)

	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go runServer(sigChan, exitChan)

	<-exitChan
}
