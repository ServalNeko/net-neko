package server

import (
	"context"
	"fmt"
	"net"
	"net-neko/input"
	"syscall"
	"time"
)

type UDP struct {
	net.UDPAddr
	input input.Input
}

func NewUDP(addr net.UDPAddr, input input.Input) *UDP {
	return &UDP{addr, input}
}

func (u *UDP) Serve() error {

	listenConfig := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if err == nil {
					return
				}

			})
		},
	}
	conn, err := listenConfig.ListenPacket(context.Background(), "udp", u.UDPAddr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		_, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		dialer := net.Dialer{
			LocalAddr: conn.LocalAddr(),
			Control: func(network, address string, c syscall.RawConn) error {
				return c.Control(func(fd uintptr) {
					err := syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
					if err == nil {
						return
					}
				})
			},
		}

		cconn, err := dialer.Dial(addr.Network(), addr.String())
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		//sub := u.input.SubScribe()
		go func(conn net.Conn) {
			defer conn.Close()
			defer print("close read")
			for {
				conn.SetDeadline(time.Now().Add(30 * time.Second))
				buf := make([]byte, 1024)
				_, err = conn.Read(buf)
				if err != nil {

					return
				}
				print(string(buf))
			}

		}(cconn)
		/*
			go func(conn net.Conn) {
				defer conn.Close()
				defer u.input.CloseSub(sub)
				defer println("close write")
				for {
					msg, ok := <-*sub
					if !ok {
						return
					}
					_, err := conn.Write([]byte(msg))
					if err != nil {
						return
					}
				}

			}(cconn)
		*/

	}
}
