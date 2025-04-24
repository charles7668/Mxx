package graceful

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var BackgroundContext context.Context

func InitContext() {
	BackgroundContext, _ = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
