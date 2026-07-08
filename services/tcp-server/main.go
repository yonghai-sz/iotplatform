package main

import (
	"os"

	"iotplatform/services/tcp-server/internal/server"
)

type serverConf struct {
	Name string    `yaml:"name"`
	Port *servPort `yaml:"port"`
}

type servPort struct {
	TCPPort string `yaml:"tcpPort"`
}

type Config struct {
	Server *serverConf `yaml:"server"`
}

func main() {

	cfg := &Config{}

	serv := server.NewTCPServer()
	addr := ":" + cfg.Server.Port.TCPPort
	err := serv.Listen(addr)
	if err != nil {
		os.Exit(1)
	}
	go serv.Start()

	select {}
}
