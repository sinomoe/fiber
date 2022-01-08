package job

import (
	rpc2 "github.com/sinomoe/fiber/pkg/comet/rpc"
	"github.com/sinomoe/fiber/pkg/queue"
	"log"
	"net/rpc"
)

type Job struct {
	cli *rpc.Client
	q   queue.Interface

	stop chan struct{}
}

func NewJob(cli *rpc.Client, q queue.Interface) *Job {
	return &Job{
		cli:  cli,
		q:    q,
		stop: make(chan struct{}),
	}
}

func (j *Job) Spin() {
	for {
		select {
		case <-j.stop:
			j.q.Shutdown()
			return
		case msg := <-j.q.Consume():
			var (
				resp rpc2.SendResponse
				err  error
			)
			if err = j.cli.Call("CometService.Send", rpc2.SendRequest{
				From:    msg.From,
				To:      msg.To,
				Message: msg.Message,
			}, &resp); err != nil {
				log.Println(err)
			}
		}
	}
}
