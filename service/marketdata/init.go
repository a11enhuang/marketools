package marketdata

import (
	"context"
	"crypto/tls"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/network/standard"
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

func AsyncStart(ctx context.Context) {
	startSyncLastPrice()
	go registerBuyStocks(ctx)
}
