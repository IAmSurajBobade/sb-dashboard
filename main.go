package main

import (
	"embed"

	"github.com/IAmSurajBobade/sb-dashboard/cmd/server"
)

//go:embed templates/*
var content embed.FS

func main() {
	// Start the server
	server.Start(content)
}
