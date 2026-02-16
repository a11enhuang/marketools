package main

import (
	"context"

	"com.reopenai/marketool/schedules"
	"com.reopenai/marketool/service/marketdata"

	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()
	ctx, cancel := context.WithCancel(context.Background())

	schedules.AsyncStart(ctx)
	marketdata.AsyncStart(ctx)

	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) { cancel() })
	h.Spin()
}
