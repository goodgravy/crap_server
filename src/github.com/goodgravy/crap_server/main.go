package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"
)

const (
	connHost                 = "localhost"
	connType                 = "tcp"
	defaultPort              = 10000
	defaultMaxPreReadDelay   = 30
	defaultMaxPreWriteDelay  = 30
	defaultSuccessPercentage = 80
)

func main() {
	port, maxPreReadDelay, maxPreWriteDelay, successPercentage := parseFlags()
	log.Printf("configuration: maxPreReadDelay(%d) maxPreWriteDelay(%d) successPercentage(%d)\n", maxPreReadDelay, maxPreWriteDelay, successPercentage)
	rand.Seed(time.Now().UnixNano())

	listener := listen(port)
	defer listener.Close()

	for {
		log.Println("blocking to accept")
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error accepting: %s\n", err.Error())
		}

		go handleClient(conn, maxPreReadDelay, maxPreWriteDelay, successPercentage)
	}
}

func parseFlags() (port, maxPreReadDelay, maxPreWriteDelay, successPercentage int) {
	flag.IntVar(&port, "port", defaultPort, "port to listen on")
	flag.IntVar(&maxPreReadDelay, "maxPreReadDelay", 30, "maximum time in seconds to wait before reading from the socket")
	flag.IntVar(&maxPreWriteDelay, "maxPreWriteDelay", 30, "maximum time in seconds to wait before writing to the socket")
	flag.IntVar(&successPercentage, "successPercentage", 80, "percentage of connections to successfully handle")
	flag.Parse()
	return
}

func listen(port int) (listener net.Listener) {
	address := fmt.Sprintf("%s:%d", connHost, port)
	listener, err := net.Listen(connType, address)
	if err != nil {
		log.Fatalf("error listening: %s\n", err.Error())
	}
	log.Printf("listening on %s\n", address)
	return
}

func waitForUpTo(message string, maxTime int) {
	sleepTime := time.Duration(rand.Intn(maxTime)) * time.Second
	log.Printf("%s, sleeping for %s\n", message, sleepTime)
	time.Sleep(sleepTime)
}

func handleClient(conn net.Conn, maxPreReadDelay, maxPreWriteDelay, successPercentage int) {
	if rand.Intn(100) > successPercentage {
		log.Println("not going to respond to client")
		return // without closing connection
	}
	defer conn.Close()

	waitForUpTo("client connected", maxPreReadDelay)
	log.Println("about to handle")

	buf := make([]byte, 1024)
	if _, err := conn.Read(buf); err != nil {
		log.Printf("Error reading: %s\n", err.Error())
	}

	waitForUpTo("request received", maxPreWriteDelay)
	conn.Write([]byte(fmt.Sprintf("replying to: %s", buf)))
}
