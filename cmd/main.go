package main

import (
	"github.com/negarciacamilo/tcp_chat/internal/logger"
	"github.com/negarciacamilo/tcp_chat/internal/server"
)

func main() {
	logger.New()
	server.StartServer()
}
