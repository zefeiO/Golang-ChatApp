package main

import (
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


func runServer(sig chan os.Signal, exit chan int) {
	broadcast := make(chan Message)
	session := Session{}
	router := mux.NewRouter()

	router.HandleFunc("/join", func(w http.ResponseWriter, r *http.Request) {
		makeSocket(w, r, &session, broadcast)
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
	go session.run(broadcast)

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
