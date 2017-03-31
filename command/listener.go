package command

import (
	"bufio"
	"fmt"
	"net"
)

type handler interface {
	Handle([]byte)
}

type Receiver struct {
	listen   string
	handlers []handler
	errors   chan error
}

func NewReceiver(listen string) *Receiver {
	return &Receiver{listen, make([]handler, 1), nil}
}

func (r *Receiver) AddHandler(h handler) {
	r.handlers = append(r.handlers, h)
}

func (r *Receiver) SetErrorChannel(c chan error) {
	r.errors = c
}

func (r *Receiver) Receive() error {
	l, err := net.Listen("tcp", r.listen)
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go r.handleRequest(conn)
	}
}

func (r *Receiver) handleRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadBytes('\n')
		if err != nil {
			r.handleError(err)
			return
		}

		go r.Dispatch(msg)
	}
}

func (r *Receiver) Dispatch(msg []byte) {
	for _, l := range r.handlers {
		l.Handle(msg)
	}
}

func (r *Receiver) handleError(err error) {
	if r.errors != nil {
		r.errors <- err
		return
	}

	fmt.Printf("Error handling request: %v \n", err)
}
