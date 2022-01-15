package main

import (
	"fmt"
	nrpc "net/rpc"

	config2 "github.com/sinomoe/fiber/internal/config"

	"github.com/sinomoe/fiber/internal/comet"
	"github.com/sinomoe/fiber/internal/comet/rpc"
)

func main() {
	cfg := new(config2.Comet)
	err := config2.Load("configs/comet.yaml", cfg)
	if err != nil {
		panic(err)
	}
	c := comet.NewComet(fmt.Sprintf(":%d", cfg.WebsocketPort))
	s := rpc.NewCometRpcServer(cfg.Rpc.Network, fmt.Sprintf(":%d", cfg.Rpc.Port))
	nrpc.Register(rpc.NewCometService(c))
	s.Run()
	c.Spin()
}
