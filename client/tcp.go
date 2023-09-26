package client

import (
	"io"
	"net"
	"net-neko/input"
	"sync"
)

type TCP struct {
	net.TCPAddr
	input input.Input
}

func NewTCP(addr net.TCPAddr, input input.Input) *TCP {
	return &TCP{addr, input}
}

func (t *TCP) Dial() error {
	conn, err := net.DialTCP("tcp", nil, &t.TCPAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	cchan := make(chan interface{})

	wg := sync.WaitGroup{}
	wg.Add(2)

	sub := t.input.SubScribe()
	go func() {
		defer wg.Done()
		defer t.input.CloseSub(sub)
		for {
			select {
			case <-cchan:
				return
			case msg, ok := <-*sub:
				if !ok {
					return
				}
				_, err := conn.Write([]byte(msg))
				if err != nil {
					return
				}

			}
		}

	}()

	go func() {
		defer wg.Done()
		defer close(cchan)
		for {
			buf := make([]byte, 1024)
			_, err := conn.Read(buf)
			if err == io.EOF {
				return
			}
			if err != nil {
				return
			}

			print(string(buf))
		}
	}()

	wg.Wait()
	return nil
}
