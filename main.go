package main

import (
	"flag"

	"com.reopenai/marketool/tencent"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")

	h := server.Default()
	go tencent.Run()
	h.Spin()
}
