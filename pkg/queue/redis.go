package queue

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/sinomoe/fiber/pkg/base"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Redis struct {
	ctx    context.Context
	stop   context.CancelFunc
	rdb    *redis.Client
	group  string
	stream string
	id     string

	produceCh chan base.Message
	consumeCh chan base.Message
}

func NewRedis(address, password, stream, group string, db int) *Redis {
	ctx, cancel := context.WithCancel(context.Background())
	r := &Redis{
		ctx:  ctx,
		stop: cancel,
		rdb: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password, // no password set
			DB:       db,       // use default DB
		}),
		group:     group,
		stream:    stream,
		id:        "consumer-" + strconv.Itoa(rand.Int()),
		produceCh: make(chan base.Message, 16),
		consumeCh: make(chan base.Message, 16),
	}
	return r
}

func (r *Redis) StartConsumer() {
	go r.run()
}

func (r *Redis) run() {
	for {
		select {
		case <-r.ctx.Done():
			return
		default:
			ss, err := r.rdb.XReadGroup(r.ctx, &redis.XReadGroupArgs{
				Group:    r.group,
				Consumer: r.id,
				Streams:  []string{r.stream, ">"},
				Count:    1,
				Block:    0,
			}).Result()
			if err != nil {
				continue
			}
			for _, ms := range ss {
				for _, m := range ms.Messages {
					for _, v := range m.Values {
						str := v.(string)
						var msg base.Message
						if err = json.Unmarshal([]byte(str), &msg); err != nil {
							break
						}
						r.consumeCh <- msg
					}
				}
			}
		}
	}
}

func (r *Redis) Produce(message base.Message) (err error) {
	var bs []byte
	if bs, err = json.Marshal(message); err != nil {
		return err
	}
	if _, err = r.rdb.XAdd(r.ctx, &redis.XAddArgs{
		Stream: r.stream,
		ID:     "*",
		Values: []string{"f", string(bs)},
	}).Result(); err != nil {
		return err
	}
	return nil
}

func (r *Redis) Consume() <-chan base.Message {
	return r.consumeCh
}

func (r *Redis) Shutdown() {
	r.stop()
}
