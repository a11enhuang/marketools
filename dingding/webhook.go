package dingding

import (
	"context"
	"crypto/tls"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/spf13/viper"
)

var httpClient *client.Client

func init() {
	clientCfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	c, _ := client.NewClient(
		client.WithTLSConfig(clientCfg),
		client.WithKeepAlive(true),
		client.WithDialer(standard.NewDialer()),
	)
	httpClient = c
}

func Send(content string) {
	msg := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]any{
			"title": "2026 买入股票推荐",
			"text":  content,
		},
	}

	url := viper.GetString("WEBHOOK_DINGDING")

	buff, err := sonic.Marshal(msg)
	if err != nil {
		return
	}

	req, res := &protocol.Request{}, &protocol.Response{}
	req.SetBody(buff)
	req.SetMethod("POST")
	req.SetRequestURI(url)

	_ = httpClient.Do(context.Background(), req, res)
}
