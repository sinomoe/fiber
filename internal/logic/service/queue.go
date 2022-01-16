package service

import (
	"context"
	"sync"

	"github.com/sinomoe/fiber/internal/logic/dto"
	"github.com/sinomoe/fiber/pkg/dto/base"
	"github.com/sinomoe/fiber/pkg/queue"
)

type Queue struct {
	q queue.Interface
}

var (
	queueSvc  *Queue
	queueOnce sync.Once
)

func InitQueue(q queue.Interface) {
	queueOnce.Do(func() {
		queueSvc = &Queue{
			q: q,
		}
	})
}

func GetQueue() *Queue {
	return queueSvc
}

func (q *Queue) SendMessage(ctx context.Context, from string, req dto.SendMessageRequest) (resp dto.SendMessageResponse, err error) {
	err = q.q.Produce(base.Message{
		From:    from,
		To:      req.To,
		Message: req.Message,
	})
	return
}
