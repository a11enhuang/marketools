package main

import (
	"context"

	"com.reopenai/marketool/tencent"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()
	ctx, cancel := context.WithCancel(context.Background())
	go tencent.Run(ctx)
	h.OnShutdown = append(h.OnShutdown, func(ctx context.Context) { cancel() })
	h.Spin()
}
