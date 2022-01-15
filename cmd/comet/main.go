package main

import (
	"github.com/sinomoe/fiber/internal/comet"
	"github.com/sinomoe/fiber/internal/comet/rpc"
	nrpc "net/rpc"
)

func main() {
	c := comet.NewComet(":8879")
	s := rpc.NewCometRpcServer("tcp", ":8869")
	nrpc.Register(rpc.NewCometService(c))
	s.Run()
	c.Spin()
}
