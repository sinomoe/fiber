package queue

import (
	"github.com/sinomoe/fiber/pkg/dto/base"
)

type Interface interface {
	Produce(message base.Message) error
	Consume() <-chan base.Message
	Shutdown()
}
