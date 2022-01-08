package rpc

import (
	"github.com/sinomoe/fiber/pkg/base"
	comet2 "github.com/sinomoe/fiber/pkg/comet"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type SendRequest struct {
	From    string
	To      string
	Message string
}

type SendResponse struct{}

type CometService struct {
	comet *comet2.Comet
}

func NewCometService(comet *comet2.Comet) CometService {
	return CometService{
		comet: comet,
	}
}

func (s CometService) Send(req SendRequest, resp *SendResponse) (err error) {
	var cli *comet2.Client
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

func NewCometRpcServer(network, address string) *CometRpcServer {
	return &CometRpcServer{
		network: network,
		address: address,
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
