package main

import (
	"github.com/dyslexic-charactersheets/translator/src/go/config"
	"github.com/dyslexic-charactersheets/translator/src/go/server"
)

func main() {
	server.RunTranslator(config.Config.Server.Host(), config.Config.Debug)
}
