package main

// #cgo CFLAGS: -g -Wall
import "C"

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
)

var (
	remoteMathServerProd = url.URL{Scheme: "wss", Host: "remote-math.mrmelon54.com", Path: "/"}
	remoteMathServerDev  = url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
)

var isDev uint32
var isEditor uint32
var alreadyRunning uint32

var upgrader = websocket.Upgrader{}

func main() {}

//export RemoteMathIsEditor
func RemoteMathIsEditor() {
	atomic.SwapUint32(&isDev, 1)
	atomic.SwapUint32(&isEditor, 1)
}

//export RemoteMathNotEditor
func RemoteMathNotEditor() {
	atomic.SwapUint32(&isEditor, 0)
}

//export RemoteMathInterfaceEntry
func RemoteMathInterfaceEntry() {
	if atomic.SwapUint32(&alreadyRunning, 1) == 0 {
		fmt.Println("[RemoteMathInterfaceEntry] Launching...")
		go internalGoRunner()
	} else {
		fmt.Println("[RemoteMathInterfaceEntry] Already running, ignoring this call...")
	}
}

func getUsableParams() url.URL {
	a := remoteMathServerProd
	if atomic.LoadUint32(&isEditor) == 1 {
		a = remoteMathServerDev
	}
	return a
}

func internalGoRunner() {
	var logFile io.Writer = os.Stdout
	listenAddr := ":8164"
	if atomic.LoadUint32(&isDev) == 1 {
		fmt.Println("[RemoteMathInterfaceEntry] Enabling development mode...")
		listenAddr = ":8165"
		create, err := os.Create("remote-math-interface.log")
		if err != nil {
			panic(err)
		}
		logFile = create
	}
	logger := log.New(logFile, "[RemoteMathInterface] ", log.LstdFlags)
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

		a := getUsableParams()

		logger.Printf("Connecting to %s\n", a.String())
		c2, _, err := websocket.DefaultDialer.Dial(a.String(), nil)
		if err != nil {
			logger.Printf("dial: %s\n", err)
			return
		}
		defer c2.Close()

		done := make(chan struct{}, 2)
		go forwardWs(done, c, c2)
		go forwardWs(done, c2, c)
		<-done
		logger.Printf("Closing connection\n")
	})
	logger.Fatal(http.ListenAndServe(listenAddr, nil))
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
