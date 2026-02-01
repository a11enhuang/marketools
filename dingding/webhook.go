package dingding

import (
	"context"
	"crypto/tls"
	"log"

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
	log.Println("[钉钉]准备发送webhook推送.url=", url)

	buff, err := sonic.Marshal(msg)
	if err != nil {
		log.Println("[钉钉]序列化参数出错.err=", err.Error())
		return
	}

	req, res := &protocol.Request{}, &protocol.Response{}
	req.SetBody(buff)
	req.SetMethod("POST")
	req.SetRequestURI(url)

	err = httpClient.Do(context.Background(), req, res)
	if err != nil {
		log.Println("[钉钉]发送webhook请求出错.err=", err.Error())
	}
}
