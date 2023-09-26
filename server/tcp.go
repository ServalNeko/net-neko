package server

import (
	"net"
	"net-neko/input"
)

type TCP struct {
	net.TCPAddr
	input input.Input
}

func NewTCP(addr net.TCPAddr, input input.Input) *TCP {
	return &TCP{addr, input}
}

func (t *TCP) Serve() error {
	listener, err := net.ListenTCP("tcp", &t.TCPAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer listener.Close()

		sub := t.input.SubScribe()
		go func() {
			for {
				msg, ok := <-*sub

				if !ok {
					t.input.CloseSub(sub)
					return
				}

				_, err := conn.Write([]byte(msg))
				if err != nil {
					return
				}
			}

		}()

		go func() {
			defer conn.Close()
			defer t.input.CloseSub(sub)

			for {
				data := make([]byte, 1024)
				_, err := conn.Read(data)
				if err != nil {
					return
				}
				print(string(data))
			}

		}()

	}
}
