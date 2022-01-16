package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/sinomoe/fiber/internal/config"
	"github.com/sinomoe/fiber/pkg/dto/base"

	"github.com/sinomoe/fiber/internal/comet"
	dto "github.com/sinomoe/fiber/pkg/dto/comet"
)

type CometService struct {
	comet *comet.Comet
}

func NewCometService(comet *comet.Comet) CometService {
	return CometService{
		comet: comet,
	}
}

func (s CometService) Send(req dto.SendRequest, resp *dto.SendResponse) (err error) {
	var cli *comet.Client
	if cli, err = s.comet.GetClient(req.To); err != nil {
		return err
	}
	if len(cli.Send) == cap(cli.Send) {
		return nil
	}
	cli.Send <- base.Message{
		From:    req.From,
		To:      req.To,
		Message: req.Message,
	}
	return nil
}

type CometRpcServer struct {
	l                net.Listener
	network, address string
}

func NewCometRpcServer(cfg *config.Comet) *CometRpcServer {
	return &CometRpcServer{
		network: cfg.Rpc.Network,
		address: fmt.Sprintf(":%d", cfg.Rpc.Port),
	}
}

func (r *CometRpcServer) Run() {
	l, err := net.Listen(r.network, r.address)
	if err != nil {
		log.Fatal("listen error:", err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}

			go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}()
}

func NewCometRpcClient(network, address string) (*rpc.Client, error) {
	return jsonrpc.Dial(network, address)
}
