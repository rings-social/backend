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

	Auth0Domain       string `arg:"--auth0-domain,env:AUTH0_DOMAIN,required" help:"Auth0 domain"`
	Auth0ClientId     string `arg:"--auth0-client-id,env:AUTH0_CLIENT_ID,required" help:"Auth0 client ID"`
	Auth0ClientSecret string `arg:"--auth0-client-secret,env:AUTH0_CLIENT_SECRET,required" help:"Auth0 client secret"`
}

var logger = logrus.New()

func main() {
	runMain()
}

func runMain() {
	arg.MustParse(&args)

	s, err := server.New(args.DatabaseUrl, &server.Auth0Config{
		Domain:       args.Auth0Domain,
		ClientId:     args.Auth0ClientId,
		ClientSecret: args.Auth0ClientSecret,
	},
		args.BaseUrl)
	if err != nil {
		logger.Fatal(err)
	}

	err = s.Run(args.ListenAddr)
	if err != nil {
		logger.Fatal(err)
	}
}
