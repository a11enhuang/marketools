package main

import (
	"com.reopenai/marketool/tencent"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()
	go tencent.Run()
	h.Spin()
}
