package main

// #cgo CFLAGS: -g -Wall
import "C"

import (
	"flag"
	"fmt"
	"github.com/MrMelon54/exit-reload"
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
	remoteMathServerProd = url.URL{Scheme: "wss", Host: "api.mrmelon54.com", Path: "/v1/remote-math"}
	remoteMathServerDev  = url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
)

var alreadyRunning uint32
var isRun, isDev bool
var listenAddr string
var upgrader = websocket.Upgrader{}

func main() {
	flag.BoolVar(&isRun, "run", false, "Run on startup")
	flag.BoolVar(&isDev, "dev", false, "Use developer server address")
	flag.StringVar(&listenAddr, "listen", ":8164", "Address to listen on")
	flag.Parse()
	if isRun {
		RemoteMathInterfaceEntry()
		exit_reload.ExitReload("RemoteMathInterface", func() {}, func() {})
	}
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
	if isDev {
		a = remoteMathServerDev
	}
	return a
}

func internalGoRunner() {
	var logFile io.Writer
	create, err := os.Create("remote-math-interface.log")
	if err != nil {
		logFile = os.Stdout
	} else {
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
	logger.Printf("Listening on '%s'\n", listenAddr)
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
