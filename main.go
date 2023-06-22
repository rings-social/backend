package main

import (
	server "backend/pkg"
	"github.com/alexflint/go-arg"
	"github.com/sirupsen/logrus"
)

var args struct {
	Debug       bool   `arg:"-D,--debug,env:RINGS_DEBUG" help:"enable debug mode"`
	ListenAddr  string `arg:"-l,--listen,env:RINGS_LISTEN_ADDR" default:"127.0.0.1:8081" help:"address to listen on"`
	DatabaseUrl string `arg:"--database-url,env:DATABASE_URL,required" help:"Database URL"`
	BaseUrl     string `arg:"--base-url,env:RINGS_BASE_URL,required" help:"Base URL for the main website"`
}

var logger = logrus.New()

func main() {
	runMain()
}

func runMain() {
	arg.MustParse(&args)

	s, err := server.New(args.DatabaseUrl, args.BaseUrl)
	if err != nil {
		logger.Fatal(err)
	}

	err = s.Run(args.ListenAddr)
	if err != nil {
		logger.Fatal(err)
	}
}
