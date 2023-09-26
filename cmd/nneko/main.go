package main

import (
	"fmt"
	"log"
	"net"
	"net-neko/client"
	"net-neko/input"
	"net-neko/server"
	"os"

	"github.com/urfave/cli/v2"
)

var commands = []*cli.Command{
	{
		Name:    "lsiten",
		Aliases: []string{"l"},
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "udp", Aliases: []string{"u"}, Usage: "listen udp port"},
		},
		Action: cmdListen,
	},
}

func cmdListen(cctx *cli.Context) error {

	args := cctx.Args()
	if args.Len() < 1 {
		log.Fatal("eeeee")
	}
	port := args.Get(0)

	isUdp := cctx.Bool("u")
	var in input.Input
	in = input.NewStdin()

	if !isUdp {
		ipAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", port))
		if err != nil {
			log.Fatal(err)
		}

		s := server.NewTCP(*ipAddr, in)
		err = s.Serve()
		if err != nil {
			log.Fatal(err)
		}

	} else {
		ipAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", port))
		if err != nil {
			log.Fatal(err)
		}
		s := server.NewUDP(*ipAddr, in)
		err = s.Serve()
		if err != nil {
			log.Fatal(err)
		}

	}

	return nil
}

func cmdClient(cctx *cli.Context) error {
	args := cctx.Args()
	if args.Len() < 2 {
		log.Fatal("eeeee")
	}

	ipAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", args.Get(0), args.Get(1)))
	if err != nil {
		log.Fatal(err)
	}

	var in input.Input

	in = input.NewStdin()

	c := client.NewTCP(*ipAddr, in)
	err = c.Dial()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Action = cmdClient
	app.Flags = []cli.Flag{
		&cli.BoolFlag{Name: "udp", Aliases: []string{"u"}, Usage: "UDP Client flag"},
	}

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}
