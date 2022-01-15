package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sinomoe/fiber/internal/logic"
	"github.com/sinomoe/fiber/pkg/config"
)

func main() {
	cfg, err := config.LoadLogic("configs/logic.yaml")
	if err != nil {
		panic(err)
	}
	svr := logic.NewServer(cfg)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	<-sig
	log.Println("shutting down gracefully")
	svr.Shutdown()

	<-time.After(time.Second)
}
