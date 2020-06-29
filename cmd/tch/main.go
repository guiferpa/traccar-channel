package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
)

var (
	id     string
	port   string
	hostip string
)

func scanStdin(msgCh chan *bytes.Buffer) {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			panic(err)
		}

		msgCh <- bytes.NewBuffer(line)
	}
}

func scanTCPServer(conn net.Conn, msgFromServerCh chan *bytes.Buffer) {
}

func main() {
	flag.StringVar(&id, "id", "", "id or IMEI")
	flag.StringVar(&hostip, "ip", "127.0.0.1", "set a custom host ip")
	flag.StringVar(&port, "port", "5055", "set a custom port")

	flag.Parse()

	if id == "" {
		id = uuid.New().String()
	}

	addr := fmt.Sprintf("%s:%s", hostip, port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	msgCh := make(chan *bytes.Buffer)
	go scanStdin(msgCh)

	msgFromServerCh := make(chan *bytes.Buffer)

	reader := bufio.NewReader(conn)
	go func() {
		for {
			bs, err := reader.ReadBytes('\n')
			if err != nil {
				log.Println("err from server:", err)
			}

			msgFromServerCh <- bytes.NewBuffer(bs)
		}
	}()

	log.Printf("Connection created for id: %s, addr: %s\n", id, conn.LocalAddr().String())

	for {
		select {
		case msg := <-msgFromServerCh:
			log.Print(msg)

		case msg := <-msgCh:
			httpRequestFormat := `GET /?%s HTTP/1.1
Host: %s
Connection: keep-alive
User-agent: Traccar Channel

			`

			commandsScanned := fmt.Sprintf("id=%s %s", id, msg.String())
			log.Printf("Commands scanned: %s\n", commandsScanned)

			query := strings.Join(strings.Split(commandsScanned, " "), "&")
			httpRequest := fmt.Sprintf(httpRequestFormat, query, addr)

			message := bytes.NewBufferString(httpRequest)
			if _, err := conn.Write(message.Bytes()); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
