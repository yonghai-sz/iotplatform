package main

import (
	"flag"
	"fmt"

	"iotplatform/services/tcp-server/internal/config"
	"iotplatform/services/tcp-server/internal/server"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
)

var configFile = flag.String("f", "etc/tcp-server.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.MustSetUp()

	serv := server.NewTCPServer()
	if err := serv.Listen(c.ListenOn); err != nil {
		logx.Must(err)
	}

	proc.AddShutdownListener(func() { serv.Stop() })

	go serv.Start()
	fmt.Printf("Starting tcp server at %s...\n", c.ListenOn)

	<-proc.Done()
}
