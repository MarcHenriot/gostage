package main

import (
	"github.com/MarcHenriot/gostage/server"
)

func main() {
	gostageServer := server.New()
	gostageServer.Init()
	gostageServer.Run()
}