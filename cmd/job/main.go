package main

import (
	"github.com/sinomoe/fiber/internal/comet/rpc"
	"github.com/sinomoe/fiber/internal/job"
	"github.com/sinomoe/fiber/pkg/queue"
	"log"
	rpc2 "net/rpc"
	"time"
)

const maxRetry = 10

func main() {
	var (
		cli   *rpc2.Client
		err   error
		retry int
	)
	for retry = 0; retry < maxRetry; retry++ {
		if cli, err = rpc.NewCometRpcClient("tcp", ":8869"); err == nil {
			break
		}
		log.Printf("retry %d times, waiting for comet\n", retry+1)
		time.Sleep((1 << retry) * time.Second)
	}
	if retry == maxRetry {
		log.Fatalln("exceeded max retry counts")
	}

	q := queue.NewRedis("127.0.0.1:6379", "", "mystream", "g1", 0)
	q.StartConsumer()
	j := job.NewJob(cli, q)
	j.Spin()
}
