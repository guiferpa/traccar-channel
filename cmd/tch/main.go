package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

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

func dateToString(n time.Time) (string, string, string) {
	year, month, day := n.Date()

	nyear := fmt.Sprintf("%v", year)

	nmonth := fmt.Sprintf("%v", int(month))
	if int(month) < 10 {
		nmonth = fmt.Sprintf("0%v", int(month))
	}

	nday := fmt.Sprintf("%v", day)
	if day < 10 {
		nday = fmt.Sprintf("0%v", day)
	}

	return nyear, nmonth, nday
}

func clockToString(n time.Time) (string, string, string) {
	hour, minute, second := n.UTC().Clock()

	nhour := fmt.Sprintf("%v", hour)
	if hour < 10 {
		nhour = fmt.Sprintf("0%v", hour)
	}

	nminute := fmt.Sprintf("%v", minute)
	if minute < 10 {
		nminute = fmt.Sprintf("0%v", minute)
	}

	nsecond := fmt.Sprintf("%v", second)
	if second < 10 {
		nsecond = fmt.Sprintf("0%v", second)
	}

	return nhour, nminute, nsecond
}

func main() {
	flag.StringVar(&id, "id", "", "id or IMEI")
	flag.StringVar(&hostip, "ip", "127.0.0.1", "set a custom host ip")
	flag.StringVar(&port, "port", "5011", "set a custom port")

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
			bs, err := reader.ReadBytes('\r')
			if err != nil {
				panic(err)
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
			commandsScanned := fmt.Sprintf("id=%s %s", id, msg.String())
			log.Printf("Commands scanned: %s\n", commandsScanned)

			n := time.Now()
			year, month, day := dateToString(n)
			hour, minute, second := clockToString(n)

			latitude := 0
			longitude := 0
			speed := 0
			course := 0
			satellites := 0
			isValid := 0
			odometer := 0
			power := 1
			status := 0
			index := 0

			tmpl := "SA200STT;%v;02;%s%s%s;%s:%s:%s;%v;%v;%v;%v;%v;%v;%v;%v;0;%v;%v;0;\r"
			message := fmt.Sprintf(tmpl, id, year, month, day, hour, minute, second, latitude, longitude, speed, course, satellites, isValid, odometer, power, status, index)
			buffer := bytes.NewBufferString(message)

			if _, err := conn.Write(buffer.Bytes()); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
