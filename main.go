package main

// #cgo CFLAGS: -g -Wall
import "C"

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
)

var (
	remoteMathServerProd = url.URL{Scheme: "wss", Host: "api.mrmelon54.com", Path: "/v1/remote-math"}
	remoteMathServerDev  = url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
)

var isEditor uint32
var alreadyRunning uint32

var upgrader = websocket.Upgrader{}

func main() {
	internalGoRunner(remoteMathServerDev)
}

//export RemoteMathIsEditor
func RemoteMathIsEditor() {
	atomic.SwapUint32(&isEditor, 1)
}

//export RemoteMathInterfaceEntry
func RemoteMathInterfaceEntry() {
	if atomic.SwapUint32(&alreadyRunning, 1) == 0 {
		fmt.Println("[RemoteMathInterfaceEntry] Launching...")
		a := remoteMathServerProd
		if atomic.LoadUint32(&isEditor) == 1 {
			a = remoteMathServerDev
			fmt.Println("[RemoteMathInterfaceEntry] Enabling development mode...")
		}
		go internalGoRunner(a)
	} else {
		fmt.Println("[RemoteMathInterfaceEntry] Already running, ignoring this call...")
	}
}

func internalGoRunner(remoteMathServer url.URL) {
	logger := log.New(os.Stdout, "[RemoteMathInterface] ", log.LstdFlags)
	logger.Println("Starting Websocket Reverse Proxy")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		c, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			logger.Println("upgrade:", err)
			return
		}
		defer c.Close()

		logger.Printf("connecting to %s\n", remoteMathServer)
		c2, _, err := websocket.DefaultDialer.Dial(remoteMathServer.String(), nil)
		if err != nil {
			fmt.Println("dial:", err)
			return
		}
		defer c2.Close()

		done := make(chan struct{}, 2)
		go forwardWs(done, c, c2)
		go forwardWs(done, c2, c)
		<-done
		logger.Printf("closing connection\n")
	})
	logger.Fatal(http.ListenAndServe(":8164", nil))
}

func forwardWs(done chan struct{}, cSrc, cDst *websocket.Conn) {
	defer func() {
		done <- struct{}{}
	}()
	for {
		mt, msg, err := cSrc.ReadMessage()
		if err != nil {
			return
		}
		err = cDst.WriteMessage(mt, msg)
		if err != nil {
			return
		}
	}
}
