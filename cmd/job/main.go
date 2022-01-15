package main

import (
	"fmt"
	"log"
	rpc2 "net/rpc"
	"time"

	"github.com/sinomoe/fiber/internal/comet/rpc"
	"github.com/sinomoe/fiber/internal/config"
	"github.com/sinomoe/fiber/internal/job"
	"github.com/sinomoe/fiber/pkg/queue"
)

func main() {
	var (
		err error
		cfg = new(config.Job)
	)
	if err = config.Load("configs/job.yaml", cfg); err != nil {
		panic(err)
	}
	var (
		cli   *rpc2.Client
		retry int
	)
	for retry = 0; retry < cfg.Rpc.Retry; retry++ {
		if cli, err = rpc.NewCometRpcClient(cfg.Rpc.Network, fmt.Sprintf(":%d", cfg.Rpc.Port)); err == nil {
			break
		}
		log.Printf("retry %d times, waiting for comet\n", retry+1)
		time.Sleep((1 << retry) * time.Second)
	}
	if retry == cfg.Rpc.Retry {
		log.Fatalln("exceeded max retry counts")
	}

	q := queue.NewRedis(cfg.Queue.Address, cfg.Queue.Password, cfg.Queue.Stream, cfg.Queue.Group, cfg.Queue.DB)
	q.StartConsumer()
	j := job.NewJob(cli, q)
	j.Spin()
}
