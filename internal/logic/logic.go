package logic

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sinomoe/fiber/internal/config"

	"github.com/sinomoe/fiber/internal/logic/service"

	"github.com/gin-gonic/gin"
	"github.com/sinomoe/fiber/pkg/queue"
)

type Server struct {
	q   queue.Interface
	srv *http.Server
}

func NewServer(cfg *config.Logic) *Server {
	q := queue.NewRedis(cfg.Queue.Address, cfg.Queue.Password, cfg.Queue.Stream, cfg.Queue.Group, cfg.Queue.DB)
	service.InitQueue(q)
	r := gin.Default()
	register(r)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	return &Server{
		q:   q,
		srv: srv,
	}
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
	s.q.Shutdown()
}
