package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/guiferpa/traccar-channel/engine"
	"github.com/guiferpa/traccar-channel/timeutil"
)

var (
	id      string
	hostip  string
	port    string
	timeout int64
)

func main() {
	flag.StringVar(&id, "id", "", "id or IMEI")
	flag.StringVar(&hostip, "ip", "127.0.0.1", "set a custom host ip")
	flag.StringVar(&port, "port", "5011", "set a custom port")
	flag.Int64Var(&timeout, "timeout", int64(time.Second*5), "set a custom timeout")

	flag.Parse()

	if id == "" {
		id = uuid.New().String()
	}

	eng, err := engine.New(hostip, port, time.Duration(timeout))
	if err != nil {
		panic(err)
	}

	defer eng.Conn.Close()

	t := timeutil.Now()
	year, month, day := t.DateString()
	hour, minute, second := t.UTCClockString()

	tmpl := "SA200STT;%s;02;%s%s%s;%s:%s:%s;0290;%s;%s;%s;%s;1;1;1;1;1;1;1;\r"

	latitude := "0"
	longitude := "0"
	speed := "0"
	course := "0"

	cmdc := make(chan *engine.Command)
	errc := make(chan error)

	go eng.Scan(os.Stdin, cmdc, errc)

	serverc := make(chan *bytes.Buffer)

	go eng.ListenServer(serverc, errc)

	log.Printf("Connection created for id: %s, addr: %s\n", id, eng.Conn.LocalAddr().String())

	for {
		select {
		case message := <-serverc:
			log.Print("[From server]:", message)

		case cmd := <-cmdc:
			if cmd == nil {
				continue
			}

			if cmd.Name == string(engine.Add) {
				if v, ok := cmd.Params["lat"]; ok {
					latitude = v
				}

				if v, ok := cmd.Params["lon"]; ok {
					longitude = v
				}

				if v, ok := cmd.Params["spd"]; ok {
					speed = v
				}

				if v, ok := cmd.Params["crs"]; ok {
					course = v
				}
			}

			if cmd.Name == string(engine.Commit) {
				message := fmt.Sprintf(tmpl, id, year, month, day, hour, minute, second, latitude, longitude, speed, course)
				buffer := bytes.NewBufferString(message)
				log.Println("[To server]:", buffer.String())

				if _, err := eng.Conn.Write(buffer.Bytes()); err != nil {
					errc <- err
				}
			}

		case err := <-errc:
			if err == io.EOF {
				log.Println("[Error]: EOF from Server")
				return
			}
			log.Println("ERROR:", err)
		}
	}
}
