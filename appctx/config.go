package appctx

import (
	"flag"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var env Env

func init() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
}

type Env struct {
}

func GetEnv() *Env {
	return &env
}

func (*Env) GetString(key string) string {
	return viper.GetString(key)
}

func (this *Env) SetDefault(key string, value any) *Env {
	viper.SetDefault(key, value)
	return this
}
