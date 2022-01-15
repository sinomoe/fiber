package job

import (
	"log"
	"net/rpc"

	"github.com/sinomoe/fiber/pkg/dto/comet"
	"github.com/sinomoe/fiber/pkg/queue"
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
				resp comet.SendResponse
				err  error
			)
			if err = j.cli.Call("CometService.Send", comet.SendRequest{
				From:    msg.From,
				To:      msg.To,
				Message: msg.Message,
			}, &resp); err != nil {
				log.Println(err)
			}
		}
	}
}
