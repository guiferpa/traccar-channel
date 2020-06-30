package engine

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type CommandString string

const (
	Add    CommandString = "add"
	Commit CommandString = "commit"
)

type Command struct {
	Name   string
	Params map[string]string
}

type Engine struct {
	Conn net.Conn
}

func parseCommand(s string) (*Command, error) {
	cmd := &Command{}

	sl := strings.Fields(s)
	if len(sl) < 1 {
		return cmd, nil
	}

	cmd.Name = sl[0]

	if cmd.Name == string(Commit) {
		return cmd, nil
	}

	if cmd.Name == string(Add) {
		params := sl[1:]
		if (len(params) % 2) == 1 {
			return nil, errors.New("invalid params amount")
		}

		cmd.Params = make(map[string]string, 0)
		i := 0
		j := 0
		for i < len(params) {
			j = i + 1
			cmd.Params[params[i]] = params[j]
			i = j + 1
		}

		return cmd, nil
	}

	return nil, nil
}

func (e *Engine) Scan(rd io.Reader, cmdc chan *Command, errc chan error) {
	reader := bufio.NewReader(rd)

	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			errc <- err
			continue
		}

		cmd, err := parseCommand(s)
		if err != nil {
			errc <- err
			continue
		}

		cmdc <- cmd
	}
}

func (e *Engine) ListenServer(serverc chan *bytes.Buffer, errc chan error) {
	reader := bufio.NewReader(e.Conn)

	for {
		bs, err := reader.ReadBytes('\r')
		if err != nil {
			errc <- err
		}

		serverc <- bytes.NewBuffer(bs)
	}
}

func New(hostip, port string, timeout time.Duration) (*Engine, error) {
	addr := fmt.Sprintf("%s:%s", hostip, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}

	return &Engine{Conn: conn}, nil
}
