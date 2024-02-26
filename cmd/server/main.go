package main

import (
	"context"
	"log"

	"github.com/Vladimir77715/my-tcp-pow/internal/config"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server/handlers"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server/handlers/receive_solution"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server/handlers/send_task"
	"github.com/Vladimir77715/my-tcp-pow/internal/tcp/server/handlers/wisdom_request"
)

func main() {
	cfg, err := config.ParseServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	hdlrs := map[int]handlers.Handler{
		server.RequestChallenge: send_task.New(int(cfg.MinSolutionRange), int(cfg.MaxSolutionRange)),
		server.SendSolution:     receive_solution.New(),
		server.WisdomRequest:    wisdom_request.New(),
	}

	if err = server.New(cfg, hdlrs).Start(ctx); err != nil {
		log.Fatal(err)
	}
}
