package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/gfelixc/gigapipe/cmd/service/bootstrap"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	err := bootstrap.StartService(ctx)
	if errors.Is(err, context.Canceled) {
		os.Exit(0)
	}

	if err != nil {
		println("error starting service:", err.Error())
		os.Exit(1)
	}
}
