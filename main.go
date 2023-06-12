package main

import (
	server "backend/pkg"
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
)

var args struct {
	Debug      bool   `arg:"-D,--debug,env:RINGS_DEBUG" help:"enable debug mode"`
	ListenAddr string `arg:"-l,--listen,env:RINGS_LISTEN_ADDR" default:"127.0.0.1:8081" help:"address to listen on"`
	Dsn        string `arg:"-d,--dsn,env:RINGS_DSN" help:"database connection string"`
}

var logger = logrus.New()

func main() {
	runMain()
}

func runMain() {
	arg.MustParse(&args)

	s, err := server.New(args.Dsn)
	if err != nil {
		logger.Fatal(err)
	}

	err = s.Run(args.ListenAddr)
	if err != nil {
		logger.Fatal(err)
	}
}
